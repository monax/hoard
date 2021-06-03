import * as grpc from '@grpc/grpc-js';
import {BinaryReader} from "google-protobuf";
import * as api from '../proto/api_grpc_pb';
import {
  Address,
  Ciphertext,
  GrantAndGrantSpec,
  Header,
  Plaintext,
  PlaintextAndGrantSpec,
  ReferenceAndCiphertext,
  ReferenceAndGrantSpec,
} from '../proto/api_pb';
import {Grant, OpenPGPSpec, PlaintextSpec, Spec, SymmetricSpec} from '../proto/grant_pb';
import {Ref} from '../proto/reference_pb';
import {StatInfo} from '../proto/stores_pb';
import {pipeline} from './pipeline';
import {BytesLike, BytesReadable, cancelAndDestroy, HeaderStream, Readable, ReadableLike} from './stream';
import {
  bytesReadable,
  bytesToObject,
  mapStream,
  objectToBytes,
  pullBytesFromObjects,
  pushBytesToObjects,
  readable,
  readBytes,
} from './streaming';

export {
  Plaintext,
  Ref,
  Grant,
  Address,
  Ciphertext,
  Header,
  Spec,
  PlaintextSpec,
  SymmetricSpec,
  OpenPGPSpec,
  PlaintextAndGrantSpec,
  ReferenceAndCiphertext,
  StatInfo,
};

const MiB = 1 << 20;

export class PlaintextStream implements HeaderStream<Header> {
  constructor(public readonly body: BytesReadable & { cancel?(): void }, public readonly head?: Header) {
  }

  bytes(): Promise<Buffer> {
    return readBytes(this.body);
  }

  headOnly(): Header | undefined {
    cancelAndDestroy(this.body);
    return this.head;
  }
}

export type Options = {
  sendChunkSize: number;
  receiveChunkSize: number;
};

const defaultOptions: Options = {
  sendChunkSize: 3 * MiB,
  receiveChunkSize: 0, // flush immediately
};

export class Client {
  readonly cleartext: api.CleartextClient;
  readonly encryption: api.EncryptionClient;
  readonly storage: api.StorageClient;
  readonly grant: api.GrantClient;

  private readonly options: Options;

  constructor(hoardURL: string, opts?: Partial<Options>) {
    const credentials = grpc.credentials.createInsecure();
    this.cleartext = new api.CleartextClient(hoardURL, credentials);
    this.encryption = new api.EncryptionClient(hoardURL, credentials);
    this.storage = new api.StorageClient(hoardURL, credentials);
    this.grant = new api.GrantClient(hoardURL, credentials);
    this.options = {...defaultOptions, ...opts};
  }

  // Cleartext

  put(body: BytesLike, header?: Header): Readable<Ref> {
    const stream = this.cleartext.put();
    if (header) {
      stream.write(make(Plaintext, p => p.setHead(header)))
    }
    return pipeline(pushPlaintexts(body, this.options.sendChunkSize), stream);
  }

  get(refs: ReadableLike<Ref>): Promise<PlaintextStream> {
    return pullPlaintexts(pipeline(readable(refs), this.cleartext.get()), this.options.receiveChunkSize);
  }

  // Encryption

  encrypt(body: BytesLike, header?: Header): Readable<ReferenceAndCiphertext> {
    const stream = this.encryption.encrypt();
    if (header) {
      stream.write(make(Plaintext, p => p.setHead(header)))
    }
    return pipeline(pushPlaintexts(body, this.options.sendChunkSize), stream);
  }

  decrypt(refs: ReadableLike<ReferenceAndCiphertext>): Promise<PlaintextStream> {
    return pullPlaintexts(pipeline(readable(refs), this.encryption.decrypt()), this.options.receiveChunkSize);
  }

  // Storage

  stat(ref: Ref): Promise<StatInfo> {
    return promisify((callback) =>
      this.storage.stat(
        make(Address, (a) => a.setAddress(ref.getAddress_asU8())),
        callback,
      ),
    );
  }

  push(encryptedData: BytesLike): Readable<Address> {
    return pipeline(
      bytesReadable(encryptedData),
      bytesToObject((buf) => make(Ciphertext, (ct) => ct.setEncrypteddata(buf)), this.options.sendChunkSize),
      this.storage.push(),
    );
  }

  pull(addresses: ReadableLike<Address>): BytesReadable {
    return pipeline(
      readable(addresses),
      this.storage.pull(),
      objectToBytes((ct: Ciphertext) => ct.getEncrypteddata_asU8(), this.options.receiveChunkSize),
    );
  }

  // Grant

  seal(spec: Spec, refs: ReadableLike<Ref>): Promise<Grant> {
    return promisify((callback) => {
        const stream = this.grant.seal(callback);
        stream.write(make(ReferenceAndGrantSpec, (rgs) => rgs.setGrantspec(spec)))
        return pipeline(
          readable(refs),
          mapStream((ref: Ref) => make(ReferenceAndGrantSpec, (rgs) => rgs.setReference(ref)), {objectMode: true}),
          stream,
        );
      },
    );
  }

  unseal(grt: Grant): Readable<Ref> {
    return this.grant.unseal(grt);
  }

  reseal(grt: Grant, spec: Spec): Promise<Grant> {
    return promisify((callback) =>
      this.grant.reseal(
        make(GrantAndGrantSpec, (gs) => {
          gs.setGrant(grt);
          gs.setGrantspec(spec);
        }),
        callback,
      ),
    );
  }

  putSeal(spec: Spec, body: BytesLike, header?: Header): Promise<Grant> {
    return promisify((callback) => {
      const stream = this.grant.putSeal(callback);

      stream.write(make(
        PlaintextAndGrantSpec,
        (pgs) => pgs.setGrantspec(spec),
        (pgs) => pgs.setPlaintext(make(Plaintext, (pt) => pt.setHead(header))),
      ));
      return pipeline(
        pushBytesToObjects(
          body,
          (buf) => make(PlaintextAndGrantSpec, (pgs) => pgs.setPlaintext(make(Plaintext, (pt) => pt.setBody(buf)))),
          this.options.sendChunkSize,
        ),
        stream,
      );
    });
  }

  unsealGet(grt: Grant): Promise<PlaintextStream> {
    const stream = this.grant.unsealGet(grt);
    return pullPlaintexts(stream, this.options.receiveChunkSize);
  }

  unsealDelete(grt: Grant): Readable<Address> {
    return this.grant.unsealDelete(grt);
  }
}

function pushPlaintexts(body: BytesLike, chunkSize: number): Readable<Plaintext> {
  return pushBytesToObjects(body, (buf) => make(Plaintext, (pt) => pt.setBody(buf)), chunkSize);
}

async function pullPlaintexts(stream: Readable<Plaintext>, chunkSize: number): Promise<PlaintextStream> {
  const {head, body} = await pullBytesFromObjects(
    stream,
    (pt) => pt.getBody_asU8(),
    chunkSize,
    (pt) => pt.getHead(),
  );
  return new PlaintextStream(body, head);
}

export function make<T>(cons: { new(): T }, ...init: ((t: T) => void)[]): T {
  const t = new cons();
  init.forEach((fn) => fn(t));
  return t;
}

export function promisify<T>(fn: (callback: grpc.requestCallback<T>) => unknown): Promise<T> {
  return new Promise((resolve, reject) =>
    fn((err, value) => {
      return err ? reject(err) : value ? resolve(value) : reject(new Error('No value or error returned'));
    }),
  );
}

export function serializeGrant(grant: Grant): string {
  return bufferFromGrant(grant).toString('base64');
}

export function deserializeGrant(grant: Uint8Array | string): Grant {
  if (typeof grant === 'string') {
    grant = Buffer.from(grant, 'base64');
  }
  // For reasons that are not clear to me Grant.deserializeBinary fails to pass instanceof Grant by the
  // time it gets to GRPC thus throwing an error
  const reader = new BinaryReader(grant);
  return Grant.deserializeBinaryFromReader(new Grant(), reader);
}

export function bufferFromGrant(grant: Grant): Buffer {
  return Buffer.from(grant.serializeBinary());
}

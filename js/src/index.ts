import * as grpc from '@grpc/grpc-js';
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
import { Grant, OpenPGPSpec, PlaintextSpec, Spec, SymmetricSpec } from '../proto/grant_pb';
import { Ref } from '../proto/reference_pb';
import { StatInfo } from '../proto/stores_pb';
import {
  BytesLike,
  BytesReadable,
  bytesReadable,
  bytesToObject,
  cancelAndDestroy,
  DEFAULT_CHUNK_SIZE,
  HeaderStream,
  objectToBytes,
  pipeline,
  pullBytesFromObjects,
  pushBytesToObjects,
  readable,
  Readable,
  ReadableLike,
} from './stream';

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

export type PlaintextStream = HeaderStream<Header>;

export class Client {
  readonly cleartext: api.CleartextClient;
  readonly encryption: api.EncryptionClient;
  readonly storage: api.StorageClient;
  readonly grant: api.GrantClient;

  constructor(hoardURL: string) {
    const credentials = grpc.credentials.createInsecure();
    this.cleartext = new api.CleartextClient(hoardURL, credentials);
    this.encryption = new api.EncryptionClient(hoardURL, credentials);
    this.storage = new api.StorageClient(hoardURL, credentials);
    this.grant = new api.GrantClient(hoardURL, credentials);
  }

  // Cleartext
  put(body: BytesLike, header?: Header): Readable<Ref> {
    return pipeline(pushPlaintexts(body, header), this.cleartext.put());
  }

  get(refs: ReadableLike<Ref>): Promise<PlaintextStream> {
    return pullPlaintexts(pipeline(readable(refs), this.cleartext.get()));
  }

  // Encryption

  encrypt(body: BytesLike, header?: Header): Readable<ReferenceAndCiphertext> {
    return pipeline(pushPlaintexts(body, header), this.encryption.encrypt());
  }

  decrypt(refs: ReadableLike<ReferenceAndCiphertext>): Promise<PlaintextStream> {
    return pullPlaintexts(pipeline(readable(refs), this.encryption.decrypt()));
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
      bytesToObject((buf) => make(Ciphertext, (ct) => ct.setEncrypteddata(buf))),
      this.storage.push(),
    );
  }

  pull(addresses: ReadableLike<Address>): BytesReadable {
    return pipeline(
      readable(addresses),
      this.storage.pull(),
      objectToBytes((ct: Ciphertext) => ct.getEncrypteddata_asU8()),
    );
  }

  // Grant

  seal(rgs: ReadableLike<ReferenceAndGrantSpec>): Promise<Grant> {
    return promisify((callback) => pipeline(readable(rgs), this.grant.seal(callback)));
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
      return pipeline(
        pushBytesToObjects(
          body,
          (buf) => make(PlaintextAndGrantSpec, (pgs) => pgs.setPlaintext(make(Plaintext, (pt) => pt.setBody(buf)))),
          make(
            PlaintextAndGrantSpec,
            (pgs) => pgs.setGrantspec(spec),
            (pgs) => pgs.setPlaintext(make(Plaintext, (pt) => pt.setHead(header))),
          ),
        ),
        this.grant.putSeal(callback),
      );
    });
  }

  unsealGet(grt: Grant): Promise<PlaintextStream> {
    const stream = this.grant.unsealGet(grt);
    return pullPlaintexts(stream);
  }

  unsealDelete(grt: Grant): Readable<Address> {
    return this.grant.unsealDelete(grt);
  }
}

// Read only the header, ensures the stream is cancelled
export function readHeader({ head, body }: PlaintextStream): Header | undefined {
  cancelAndDestroy(body);
  return head;
}

function pushPlaintexts(body: BytesLike, header?: Header): Readable<Plaintext> {
  return pushBytesToObjects(
    body,
    (buf) => make(Plaintext, (pt) => pt.setBody(buf)),
    make(Plaintext, (pt) => pt.setHead(header)),
  );
}

function pullPlaintexts(stream: Readable<Plaintext>): Promise<PlaintextStream> {
  return pullBytesFromObjects(
    stream,
    (pt) => pt.getBody_asU8(),
    (pt) => pt.getHead(),
  );
}

export function make<T>(cons: { new (): T }, ...init: ((t: T) => void)[]): T {
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

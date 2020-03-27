import { Readable, Transform, TransformOptions } from 'stream';
import * as api from "../proto/api_grpc_pb";
import * as grpc from "grpc";
import { Plaintext, Ciphertext, PlaintextAndGrantSpec, ReferenceAndCiphertext, Address, Header } from "../proto/api_pb";
import { Spec, PlaintextSpec, SymmetricSpec, OpenPGPSpec } from "../proto/grant_pb";
import { StatInfo } from "../proto/stores_pb";
import { Ref } from "../proto/reference_pb";
import { Grant } from "../proto/grant_pb";

export { 
  Plaintext, Ref, Grant, Address, Ciphertext, Header,
  Spec, PlaintextSpec, SymmetricSpec, OpenPGPSpec,
  PlaintextAndGrantSpec, ReferenceAndCiphertext, StatInfo, 
}

type resolver<T> = (value?: T | PromiseLike<T>) => void
type rejecter = (reason?: any) => void

const DEFAULT_CHUNK_SIZE = 2 ** 16

function extend(...args: any[]): any {
  const next = {};
  for (const obj of args) {
      for (const key in obj) {
        next[key] = obj[key];
      }
  }
  return next;
};

export type Client = api.CleartextClient & api.EncryptionClient & api.StorageClient & api.GrantClient

export function Client(hoardURL: string): Client {
  let cli: Client = {} as any;
  for (const ServiceClient of [
    api.CleartextClient,
    api.EncryptionClient,
    api.StorageClient,
    api.GrantClient,
  ]) {
    cli = extend(cli, new ServiceClient(hoardURL, grpc.credentials.createInsecure()));
  }
  return cli;
}

export function NewHeader(salt?: string | Uint8Array): Header {
  const msg = new Header();
  msg.setSalt(salt);
  return msg
}

export function NewPlaintext(data?: string | Uint8Array, head?: Header): Plaintext {
  const msg = new Plaintext();
  msg.setBody(data);
  msg.setHead(head);
  return msg
}

export function NewAddress(data?: Uint8Array): Address {
  const msg = new Address();
  msg.setAddress(data);
  return msg
}

export function NewPlaintextAndGrantSpec(pt?: Plaintext, spec?: Spec): PlaintextAndGrantSpec {
  const msg = new PlaintextAndGrantSpec();
  msg.setPlaintext(pt);
  msg.setGrantspec(spec);
  return msg;
}

export function NewPlaintextSpec(): Spec {
  let spec = new Spec();
  spec.setPlaintext(new PlaintextSpec());
  return spec;
}

export function NewSymmetricSpec(id: string): Spec {
    let symmetricSpec = new SymmetricSpec();
    symmetricSpec.setPublicid(id);
    let spec = new Spec();
    spec.setSymmetric(symmetricSpec);
    return spec;
}

export function NewOpenPGPSpec(pubKey: string): Spec {
    let openpgpSpec = new OpenPGPSpec();
    openpgpSpec.setPublickey(pubKey);
    let spec = new Spec();
    spec.setOpenpgp(openpgpSpec);
    return spec;
}

export function ChunkData(data: Buffer, chunkSize = DEFAULT_CHUNK_SIZE): Uint8Array[] {
  const chunks: Uint8Array[] = [];
  while (data.length > 0) {
    chunks.push(data.slice(0, chunkSize));
    data = data.slice(chunkSize);
  }
  return chunks;
}

function merge(left: Uint8Array, right: Uint8Array): Uint8Array {
  let data = new Uint8Array(left.length + right.length);
  data.set(left);
  data.set(right, left.length);
  return data
}

export function ReducePlaintext(all: Plaintext[]): Plaintext {
  return all.reduce((left, right) => {
    left = left || new Plaintext();
    right = right || new Plaintext();
    return NewPlaintext(
      merge(left.getBody_asU8(), right.getBody_asU8()), 
      left.getHead() || right.getHead()
    );
  }, new Plaintext())
}

export function ReduceCiphertext(all: Ciphertext[]): Ciphertext {
  return all.reduce((left, right) => {
    const ciphertext = new Ciphertext();
    left = left || new Ciphertext();
    right = right || new Ciphertext();
    ciphertext.setEncrypteddata(merge(left.getEncrypteddata_asU8(), right.getEncrypteddata_asU8()));
    return ciphertext;
  }, new Ciphertext())
}

export function ReduceReferenceAndCiphertext(all: ReferenceAndCiphertext[]): ReferenceAndCiphertext {
  return all.reduce((left, right) => {
    const refAndCT = new ReferenceAndCiphertext();
    refAndCT.setReference(left.getReference() || right.getReference());
    refAndCT.setCiphertext(ReduceCiphertext([left.getCiphertext(), right.getCiphertext()]));
    return refAndCT;
  }, new ReferenceAndCiphertext())
}

function read<T>(accum: T[], stream: Readable, resolve: resolver<T[]>, reject: rejecter): void {
  stream.on('data', (data: T) => { accum = accum.concat(data) });
  stream.on('error', (err: { code: grpc.status; }) => err.code === grpc.status.CANCELLED ? resolve(accum) : reject(err));
  stream.on('close', () => resolve(accum));
  stream.on('end', () => resolve(accum));
}

export async function Duplex<A, B>(input: A[], accum: B[], stream: grpc.ClientDuplexStream<A, B>): Promise<B[]> {
  return new Promise((resolve, reject) => {
    read(accum, stream, resolve, reject);
    input.map(value => stream.write(value));
    stream.end();
  })
}

export async function Read<T>(accum: T[], stream: grpc.ClientReadableStream<T>): Promise<T[]> {
  return new Promise((resolve, reject) => read(accum, stream, resolve, reject));
}

export async function ReadUntil<T>(accum: T[], stream: grpc.ClientReadableStream<T>, ret: (data: T[]) => boolean): Promise<T[]> {
  return new Promise((resolve, reject) => {
    stream.on('data', (data: T) => { 
      accum = accum.concat(data);
      if (ret(accum)) {
        stream.cancel();
        resolve(accum);
      }
    });
    stream.on('error', (err: { code: grpc.status; }) => err.code === grpc.status.CANCELLED ? resolve(accum) : reject(err));
    stream.on('close', () => resolve(accum));
    stream.on('end', () => resolve(accum));
  })
}

export async function ReadHeader(stream: grpc.ClientReadableStream<Plaintext>) {
  return new Promise<Header>((resolve, reject) => {
    stream.on('data', (data: Plaintext) => { 
      if (data.hasHead()) {
        stream.cancel();
        resolve(data.getHead());
      }
    });
    stream.on('error', (err: { code: grpc.status; }) => err.code === grpc.status.CANCELLED ? resolve() : reject(err));
    stream.on('close', () => reject("no header found"));
    stream.on('end', () => reject("no header found"));
  })
}

export async function Write<A, B>(input: A[], fn: (callback: grpc.requestCallback<B>) => grpc.ClientWritableStream<A>): Promise<B> {
  return new Promise((resolve, reject) => {
    let stream = fn((err, grt) => err ? reject(err) : resolve(grt));
    input.map(value => {
      stream.write(value);
    })
    stream.end();
  })
}

// Converts an object mode stream to a byte mode stream selecting a single buffer from the source using dataSelector
// and buffering bytes in chunks of at chunkSize
const BufferTransform = (bufferFromInput = (msg: any) => msg, bufferToOutput = (buf: any) => buf, 
    chunkSize = DEFAULT_CHUNK_SIZE, transformOptions: TransformOptions = {}) => {

  let buffer = Buffer.alloc(0);

  const push = (transform: Transform, buffer: Buffer) => {
    if (transform.readableLength == 0) {
      transform.push(bufferToOutput(buffer))
    }
  }

  const transform = new Transform({
    transform(msg, encoding, callback) {
      buffer = Buffer.concat([buffer, bufferFromInput(msg)])
      if (buffer.length > chunkSize) {
        push(this, buffer.slice(0, chunkSize))
        buffer = buffer.slice(chunkSize)
      }
      callback()
    },

    flush (callback) {
      push(this, buffer)
      callback()
    },
    ...transformOptions
  })

  transform.on('unpipe', src => {
    // Make sure we cancel any GRPC stream
    if (typeof src.cancel === 'function') {
      src.cancel()
    }
  })
  return transform
}

export const BytesToObjects = (bufferToObject = (buf: Uint8Array) => NewPlaintext(buf, null), chunkSize = DEFAULT_CHUNK_SIZE) =>
  BufferTransform(msg => msg, bufferToObject, chunkSize,
    { readableObjectMode: true, writableObjectMode: false })

// Converts an object mode stream to a byte mode stream selecting a single buffer from the source using dataSelector
// and buffering bytes in chunks at chunkSize
export const PlaintextToBytes = (chunkSize = DEFAULT_CHUNK_SIZE) =>
  BufferTransform((msg: Plaintext) => msg.getBody_asU8(), msg => msg, chunkSize,
    { readableObjectMode: false, writableObjectMode: true })

// Reads from the provided byteStream until a varint length-prefixed prefix of the stream can be returned.
// The stream will be destroyed once the prefix has been read so will not be totally consumed
export async function ReadLengthPrefixed(stream: Readable, byteLength: number): Promise<Buffer> {
  let buffer = Buffer.alloc(0);
  let length = 0;

  return new Promise((resolve, reject) => {
    stream.on('data', (data: Buffer) => {
      var buf = Buffer.concat([buffer, data])
      // First try to read the length prefix
      if (length === 0) {
        if (buf.length >= byteLength) {
          length = buf.readUIntBE(0, byteLength)
          buf = buf.slice(byteLength)
        }
      }
      if (length > 0 && buf.length >= length) {
        // If we have read the length prefix and it is contained within
        // our current buffer then we are done
        return resolve(buf.slice(0, length))
      }
      // Keep growing the buffer until it contains sufficient data
      buffer = buf
    });
    stream.on('error', (err: { code: grpc.status; }) => err.code === grpc.status.CANCELLED ? resolve(buffer) : reject(err));
    stream.on('close', () => resolve(buffer));
    stream.on('end', () => resolve(buffer));
  });
}
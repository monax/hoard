import * as grpc from '@grpc/grpc-js';
import { ObjectReadable, ObjectWritable } from '@grpc/grpc-js/build/src/object-stream';
import * as nstream from 'stream';
import { Stream, Transform, TransformOptions } from 'stream';

export const DEFAULT_CHUNK_SIZE = 2 ** 16;

export type Readable<Input = unknown> = ObjectReadable<Input>;
export type Writable<Output = unknown> = ObjectWritable<Output>;
export type Duplex<Input = unknown, Output = unknown> = Writable<Input> & Readable<Output>;
export type BytesReadable = Readable<Uint8Array>;
export type BytesLike = BytesReadable | Uint8Array | string | Iterable<number>;
export type ReadableLike<T> = Iterable<T> | Readable<T>;

export type HeaderStream<Header> = {
  head?: Header | undefined;
  body: BytesReadable & { cancel?(): void };
};

export function bytesReadable(bs: BytesLike): BytesReadable {
  if (bs instanceof nstream.Stream) {
    if (isReadable(bs)) {
      return bs;
    }
    throw new Error(`BytesLike '${bs}' is Stream but not readable`);
  }
  return nstream.Readable.from(Buffer.from(bs), { objectMode: false });
}

export function readable<T>(r: ReadableLike<T>): Readable<T> {
  if (isReadable(r)) {
    return r;
  }
  return nstream.Readable.from(r);
}

function isWritable(s: unknown): s is Writable {
  return (s as any)['writable'] ?? false;
}

function isReadable(s: unknown): s is Readable {
  return (s as any)['readable'] ?? false;
}

export const exitEarly: unique symbol = Symbol('exitEarly');

export class ReadError<O> extends Error {
  constructor(err: Error, public readonly consumed: O) {
    super(err.message);
    this.name = `ReadError(${err.name})`;
    this.stack = err.stack;
  }
}

// Stream reduce function, if reducer returns exitEarly then the current value of the accumulator will be returned
export function read<T, O = T[]>(
  stream: Readable<T>,
  reducer: (accum: O, data: T) => O | typeof exitEarly,
  accum: O,
): Promise<O> {
  return new Promise((resolve, reject) => {
    stream.on('data', (data: T) => {
      const reduced = reducer(accum, data);
      if (reduced === exitEarly) {
        cancelAndDestroy(stream);
        return resolve(accum);
      }
      accum = reduced;
    });
    stream.on('error', (err: Error & { code?: grpc.status }) =>
      err.code === grpc.status.CANCELLED ? resolve(accum) : reject(new ReadError(err, accum)),
    );
    stream.on('close', () => resolve(accum));
    stream.on('end', () => resolve(accum));
  });
}

export function readAll<T>(stream: Readable<T>, earlyExit?: (accum: T[], data: T) => boolean): Promise<T[]> {
  return read(
    stream,
    (accum, data) => (earlyExit && earlyExit(accum, data) ? exitEarly : accum.concat(data)),
    [] as T[],
  );
}

export function readBytes(stream: BytesLike): Promise<Uint8Array> {
  return read(bytesReadable(stream), (accum, data) => Buffer.concat([accum, data]), Buffer.alloc(0));
}

export function waitFor(stream: Stream): Promise<void> {
  return new Promise((resolve, reject) => {
    stream.on('error', (err: Error & { code?: grpc.status }) =>
      err.code === grpc.status.CANCELLED ? resolve() : reject(err),
    );
    stream.on('close', resolve);
    stream.on('end', resolve);
  });
}

export function passThrough<T>(): Duplex<T, T> {
  return new Transform({
    transform(input: T, encoding, callback) {
      callback(null, input);
    },
  });
}

export function mapStream<Input, Output>(
  fn: (i: Input) => Output,
  opts?: Partial<TransformOptions>,
): Duplex<Input, Output> {
  return new Transform({
    transform(input: Input, encoding, callback) {
      callback(null, fn(input));
    },
    ...opts,
  });
}

export function prefixStream<T>(prefix: T[], opts?: Partial<TransformOptions>): Duplex<T, T> {
  const [first] = prefix;
  return new Transform({
    readableObjectMode: isObjectModeElement(first),
    transform(input: T, encoding, callback) {
      // Push prefix before the first element and then replace with passthrough function
      prefix.forEach((p) => this.push(p));
      this.push(input);
      this._transform = (input: T, encoding, callback) => callback(null, input);
      callback();
    },
    ...opts,
  });
}

export function pushBytesToObjects<T>(
  body: BytesLike,
  bufferToOutput: (buf: Uint8Array) => T,
  header?: T,
): Readable<T> {
  const prefix: T[] = [];
  if (header) {
    prefix.push(header);
  }
  return pipeline(
    bytesReadable(body),
    bytesToObject(bufferToOutput),
    prefixStream(prefix, { writableObjectMode: true, readableObjectMode: true }),
  );
}

export async function pullBytesFromObjects<T, H = void>(
  stream: Readable<T>,
  objectToBuffer: (obj: T) => Uint8Array,
  getHeader: (obj: T) => H | undefined,
): Promise<HeaderStream<H>> {
  const first = await new Promise<T>((resolve, reject) => {
    stream.on('readable', () => {
      stream.pause();
      resolve(stream.read(1));
      stream.resume();
    });
    stream.on('error', reject);
  });
  const outputStream = objectToBytes(objectToBuffer);
  const head = getHeader(first);
  if (!head) {
    outputStream.write(first);
  }
  return {
    head,
    body: pipeline(stream, outputStream),
  };
}

export function bytesToObject<Output>(
  bufferToOutput: (buf: Uint8Array) => Output,
  chunkSize = DEFAULT_CHUNK_SIZE,
): Transform {
  return buffered<Uint8Array, Output>((buf) => buf, bufferToOutput, chunkSize, {
    readableObjectMode: true,
    writableObjectMode: false,
  });
}

export function objectToBytes<Input>(
  inputToBuffer: (input: Input) => Uint8Array,
  chunkSize = DEFAULT_CHUNK_SIZE,
): Transform {
  return buffered<Input, Uint8Array>(inputToBuffer, (buf) => buf, chunkSize, {
    readableObjectMode: false,
    writableObjectMode: true,
  });
}

// Converts an object mode stream to a byte mode stream selecting a single buffer from the source using dataSelector
// and buffering bytes in chunks of at chunkSize
function buffered<Input, Output>(
  inputToBuffer: (input: Input) => Uint8Array,
  bufferToOutput: (buf: Uint8Array) => Output,
  chunkSize = DEFAULT_CHUNK_SIZE,
  transformOptions: TransformOptions = {},
): Transform {
  let buffer = Buffer.alloc(0);

  const push = (transform: Transform, buffer: Buffer) => {
    if (transform.readableLength == 0) {
      const output = bufferToOutput(buffer);
      transform.push(output);
    }
  };

  const transform = new Transform({
    transform(msg, encoding, callback) {
      buffer = Buffer.concat([buffer, inputToBuffer(msg)]);
      if (buffer.length > chunkSize) {
        push(this, buffer.slice(0, chunkSize));
        buffer = buffer.slice(chunkSize);
      }
      callback();
    },

    flush(callback) {
      push(this, buffer);
      callback();
    },
    ...transformOptions,
  });

  transform.on('unpipe', cancelAndDestroy);

  return transform;
}

// Reads from the provided byteStream until a varint length-prefixed prefix of the stream can be returned.
// The stream will be destroyed once the prefix has been read so will not be totally consumed
export async function readLengthPrefixed(stream: Readable<Uint8Array>, byteLength: number): Promise<Buffer> {
  let buffer = Buffer.alloc(0);
  let prefixLength = 0;

  const prefix = await read(
    stream,
    (accum, data) => {
      if (accum) {
        // We have set the buffer
        return exitEarly;
      }
      let buf = Buffer.concat([buffer, data]);
      // First try to read the length prefix
      if (prefixLength === 0) {
        if (buf.length >= byteLength) {
          prefixLength = buf.readUIntBE(0, byteLength);
          // Chop off the length prefix itself
          buf = buf.slice(byteLength);
        }
      }
      if (buf.length >= prefixLength) {
        // If we have read the length prefix and it is contained within
        // our current buffer then we are done
        return buf.slice(0, prefixLength);
      }
      // Keep growing the buffer until it contains sufficient data
      buffer = buf;
      return undefined;
    },
    undefined as undefined | Buffer,
  );
  if (!prefix) {
    throw new Error(`Could not read length-prefixed prefix from stream`);
  }
  return prefix;
}

export function cancelAndDestroyOnError(stream: Readable | Writable): (err: Error | null) => void {
  return (err) => err && cancelAndDestroy(stream, err);
}

type Cancellable = Stream & { cancel: () => void };

function isCancellable(stream: Stream): stream is Cancellable {
  return typeof (stream as any)['cancel'] === 'function';
}

export function cancelAndDestroy(stream: Readable | Writable, err?: Error): void {
  if (isCancellable(stream)) {
    stream.cancel();
  }
  stream.destroy(err || undefined);
}

export type PipedStream<First extends Stream, Last extends Stream> = First extends Writable<infer Input>
  ? Last extends Readable<infer Output>
    ? Duplex<Input, Output>
    : Writable<Input>
  : Last extends Readable<infer Output>
  ? Readable<Output>
  : Stream;

export function pipeline<First extends Readable, Transforms extends Duplex[], Last extends Writable>(
  ...streams: [First, ...Transforms, Last]
): PipedStream<First, Last> {
  const first = streams[0];
  const last = streams[streams.length - 1];

  if (!isReadable(first)) {
    throw new Error(`First stream in pipe must be readable`);
  }

  if (!isWritable(last)) {
    throw new Error(`Last stream in pipe must be writable`);
  }

  const piped = streams.reduce(piper());

  // Mutually assured destruction
  const destroyer = once((err?: Error) => {
    streams.forEach((s) => cancelAndDestroy(s, err));
    piped.destroy(err);
  });

  streams.forEach((s, i) => s.on('error', destroyer));

  piped.on('error', destroyer);

  // We have to assert here since the available typings are not tight enough to capture this type
  return piped as PipedStream<First, Last>;
}

// Guess the appropriate object mode setting based on a candidate elemetn
function isObjectModeElement(t: unknown): boolean {
  return t !== null && t !== undefined && typeof t !== 'string' && !(t instanceof Uint8Array);
}

function piper(): <Left extends Stream, Right extends Stream>(left: Left, right: Right, index: number) => Right {
  return (left, right, index) => {
    if (!isReadable(left)) {
      throw new Error(`Stream to the left of pipe is not readable`);
    }
    if (!isWritable(right)) {
      throw new Error(`Stream to the right of pipe is not writable`);
    }
    const leftObjectMode = Boolean(left.readableObjectMode);
    const rightObjectMode = Boolean(right.writableObjectMode);
    if (leftObjectMode !== rightObjectMode) {
      throw new Error(
        `Streams in pipe do not agree on their object mode between streams ${index + 1} (${leftObjectMode}) and ${
          index + 2
        } (${rightObjectMode}) in pipe`,
      );
    }
    return left.pipe(right);
  };
}

function once<T extends unknown[]>(fn: (...args: [...T]) => void): typeof fn {
  let ran = false;
  return (...args) => {
    if (!ran) {
      ran = true;
      fn(...args);
    }
  };
}

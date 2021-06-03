import * as grpc from '@grpc/grpc-js';
import { ObjectDuplex } from '@grpc/grpc-js/build/src/object-stream';
import * as nstream from 'stream';
import { Stream, Transform, TransformOptions } from 'stream';
import { pipeline } from './pipeline';
import { Slice } from './slice';
import {
  BytesLike,
  BytesReadable,
  cancelAndDestroy,
  Duplex,
  HeaderStream,
  isReadableStreak,
  Readable,
  ReadableLike,
} from './stream';

const KiB = 1 << 10;

export function bytesReadable(bs: BytesLike): BytesReadable {
  if (bs instanceof nstream.Stream) {
    if (isReadableStreak(bs)) {
      return bs;
    }
    throw new Error(`BytesLike '${bs}' does not seem to be a readable stream`);
  }
  if (bs instanceof Uint8Array || typeof bs === 'string') {
    bs = Buffer.from(bs);
  }
  return nstream.Readable.from(bs, { objectMode: false });
}

export function readable<T>(r: ReadableLike<T>): Readable<T> {
  if (isReadableStreak(r)) {
    return r;
  }
  return nstream.Readable.from(r);
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
    (accum, data) => {
      if (earlyExit && earlyExit(accum, data)) {
        return exitEarly;
      }
      accum.push(data);
      return accum;
    },
    [] as T[],
  );
}

export async function readBytes(stream: BytesLike, sizeHint = KiB): Promise<Buffer> {
  const slice = await read(
    bytesReadable(stream),
    (accum, data) => accum.appendInPlace(data),
    new Slice(Buffer.allocUnsafe(sizeHint)),
  );
  return slice.buffer();
}

export function waitFor(stream: Stream): Promise<void> {
  return new Promise((resolve, reject) => {
    stream.on('error', (err: Error & { code?: grpc.status }) =>
      err.code === grpc.status.CANCELLED ? resolve() : reject(err),
    );
    // Writable
    stream.on('finish', resolve);
    // Readable
    stream.on('end', resolve);
    stream.on('close', resolve);
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

export function pushBytesToObjects<T>(
  body: BytesLike,
  bufferToOutput: (buf: Uint8Array) => T,
  chunkSize: number,
): Readable<T> {
  return pipeline(bytesReadable(body), bytesToObject(bufferToOutput, chunkSize));
}

export async function pullBytesFromObjects<T, H = void>(
  stream: Readable<T>,
  objectToBuffer: (obj: T) => Uint8Array,
  chunkSize: number,
  getHeader: (obj: T) => H | undefined,
): Promise<HeaderStream<H>> {
  const first = await new Promise<T>((resolve, reject) => {
    stream.on('readable', () => {
      stream.pause();
      const first = stream.read(1);
      // If the first frame is null this should mean that an error will follow so wait to reject that
      if (first) {
        resolve(first);
      }
      stream.resume();
    });
    stream.on('error', (err) => reject(err));
  });
  const outputStream = objectToBytes(objectToBuffer, chunkSize);
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
  chunkSize: number,
): ObjectDuplex<Uint8Array, Output> {
  return buffered<Uint8Array, Output>((buf) => buf, bufferToOutput, chunkSize, {
    readableObjectMode: true,
    writableObjectMode: false,
  });
}

export function objectToBytes<Input>(
  inputToBuffer: (input: Input) => Uint8Array,
  chunkSize: number,
): ObjectDuplex<Input, Uint8Array> {
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
  chunkSize: number,
  transformOptions: TransformOptions = {},
): ObjectDuplex<Input, Output> {
  let bufferOffset = 0;
  let buffer = Buffer.allocUnsafe(chunkSize);

  const flush = (transform: Transform): void => {
    transform.push(bufferToOutput(buffer));
    bufferOffset = 0;
    buffer = Buffer.allocUnsafe(chunkSize);
  };

  const transform = new Transform({
    transform(msg, encoding, callback): void {
      // Buffer.from should be zero copy here
      const input = Buffer.from(inputToBuffer(msg));

      if (chunkSize === 0) {
        return callback(null, input);
      }

      let written = 0;
      // Copy from the accum into the buffer until it is full
      while (written < input.length) {
        const n = input.copy(buffer, bufferOffset, written);
        bufferOffset += n;
        written += n;
        if (bufferOffset == chunkSize) {
          flush(this);
        }
      }

      // If there were enough bytes in accum to fill the buffer then flush the buffer
      callback();
    },

    flush(callback) {
      if (bufferOffset > 0) {
        buffer = buffer.slice(0, bufferOffset);
        flush(this);
      }
      callback();
    },
    ...transformOptions,
  });

  transform.on('unpipe', cancelAndDestroy);

  return transform;
}

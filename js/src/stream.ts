import {ObjectReadable, ObjectWritable} from '@grpc/grpc-js/build/src/object-stream';
import {Stream} from 'stream';

// The GRPC types use a hack in order to overlap with the native Node stream types
// it sort of works and allows us to use types annotations in various places so we
// make use of the object types for slightly stronger types.
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

type Cancellable = Stream & { cancel: () => void };

function isCancellable(stream: Stream): stream is Cancellable {
  return typeof (stream as any)['cancel'] === 'function';
}

export function isWritableStream(s: unknown): s is Writable {
  return (s as any)['writable'] !== undefined;
}

export function isReadableStreak(s: unknown): s is Readable {
  return (s as any)['readable'] !== undefined;
}

export function cancelAndDestroy(stream: Stream, err?: Error): void {
  if (isCancellable(stream)) {
    stream.cancel();
  }
  if (isReadableStreak(stream) || isWritableStream(stream)) {
    stream.destroy(err || undefined);
  }
}

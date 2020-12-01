import {Stream} from 'stream';
import {Duplexed} from './duplexed';
import {cancelAndDestroy, Duplex, isReadable, isWritable, Readable, Writable} from './stream';

// A pipeline that returns Duplex stream when it can...

export type PipedStream<First extends Stream, Last extends Stream> = First extends Writable<infer Input>
  ? Last extends Readable<infer Output>
    ? Duplex<Input, Output>
    : Writable<Input>
  : Last extends Readable<infer Output>
  ? Readable<Output>
  : Stream;

// Pipe together various streams and provided a combined stream of the strongest possible type
// i.e. a single duplex stream if possible, and handle errors and destruction all for one and one for all
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

  // The last stream is returned after being piped
  streams.reduce(piper());

  const piped: Stream = isWritable(first)
    ? isReadable(last)
      ? new Duplexed(first, last) // Duplex
      : first // Writable
    : isReadable(last)
    ? last // Readable
    : last; // Neither (but we still need to hang events off something

  // Mutually assured destruction
  const destroyer = once((err?: Error) => {
    streams.forEach((s) => cancelAndDestroy(s, err))
    cancelAndDestroy(piped)
  });

  streams.forEach((s) => s.on('error', destroyer));

  piped.on('error', destroyer)

  // We have to assert here since the available typings are not tight enough to capture this type
  return piped as PipedStream<First, Last>;
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

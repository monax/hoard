import * as nstream from 'stream';
import { Duplex, Readable, Writable } from './stream';

// Could not find anything that came sufficiently close to this, duplexify very convoluted duplexer2
// didn't work properly, neither in typescript, nothing in node core

export class Duplexed<Input, Output> extends nstream.Duplex implements Duplex<Input, Output> {
  private readonly _writable: Writable<Input>;
  private readonly _readable: Readable<Output>;

  // Wrap writable and readable as a single duplex stream that writes into writable and reads
  // from readable. If writable happens to be input to pipe and readable output then acts as a combined
  // stream (see pipeline usage)
  constructor(writable: Writable<Input>, readable: Readable<Output>, opts?: nstream.DuplexOptions) {
    super({
      ...opts,
      readableObjectMode: readable.readableObjectMode,
      writableObjectMode: writable.writableObjectMode,
      readableHighWaterMark: readable.readableHighWaterMark,
      writableHighWaterMark: writable.writableHighWaterMark,
    });

    this._writable = writable;
    this._readable = readable;

    // readable.read() should only be called when readable is paused so we pause() here so we can read()
    // later
    readable.pause();

    this.once('finish', () => writable.end());

    // handle end of stream
    writable.on('finish', () => this.end());
    readable.on('end', () => this.push(null));

    // forward errors
    writable.on('error', (err) => this.emit('error', err));
    readable.on('error', (err) => this.emit('error', err));
  }

  _destroy(error: Error | null): void {
    this._writable.destroy(error || undefined);
    this._readable.destroy(error || undefined);
  }

  _write(chunk: Input, encoding: BufferEncoding, callback: (error?: Error | null) => void): void {
    const ok = this._writable.write(chunk, encoding, () => ok && callback());
    if (!ok) {
      // Wait for drain if backpressured
      this._writable.once('drain', callback);
    }
  }

  _read(size?: number): void {
    const chunk = this._readable.read(size);
    if (chunk !== null) {
      this.push(chunk);
    } else {
      // Wait for readable event again on null
      this._readable.once('readable', () => this._read(size));
    }
  }
}

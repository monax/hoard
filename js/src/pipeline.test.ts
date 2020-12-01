import { Readable, Transform, TransformOptions, Writable } from 'stream';
import {pipeline} from "./pipeline";
import { readAll, waitFor } from './streaming';

describe('pipeline', () => {
  test('error propagation', async () => {
    const explode = Buffer.from('bang');

    const first = new Transform({
      transform(chunk, encoding, callback) {
        this.push(chunk);
        const err = Buffer.compare(explode, chunk) === 0 ? new Error(`We exploded: ${chunk.toString()}`) : undefined;
        callback(err);
      },
    });

    const second = new Transform({
      transform(chunk, encoding, callback) {
        this.push(chunk);
        callback();
      },
    });

    const firstError = new Promise((resolve, reject) => first.on('error', (err) => reject(err)));
    const secondError = new Promise((resolve, reject) => first.on('error', (err) => reject(err)));

    first.write('foo');
    first.write(explode);
    const p = pipeline(first, second);
    await expect(readAll(p)).rejects.toThrowError(explode.toString());
    // Expect error to surface

    expect(first.destroyed).toEqual(true);
    expect(second.destroyed).toEqual(true);

    // Both streams should err
    await expect(firstError).rejects.toThrowError(explode.toString());
    await expect(secondError).rejects.toThrowError(explode.toString());
  });

  test('non-duplex', async () => {
    // Test the case when the result of pipeline is not really a Duplex (and find a way to live with yourself)
    const input = ['one', 'two', 'three'];
    const output: string[] = [];
    let i = 0;

    const first = new Readable({
      read(size: number) {
        if (i >= input.length) {
          this.push(null);
          return;
        }
        this.push(input[i++]);
      },
    });

    const second = new Writable({
      write(chunk, encoding, callback) {
        output.push(chunk.toString());
        callback();
      },
    });

    await waitFor(pipeline(first, second));

    expect(first.destroyed).toEqual(true);
    expect(second.destroyed).toEqual(true);

    expect(output).toEqual(input);
  });

  test('object duplex', async () => {
    const opts: TransformOptions = {
      readableObjectMode: true,
      writableObjectMode: true,
    };

    const first = new Transform({
      transform({ foo }, encoding, callback) {
        callback(undefined, { foo: foo + 2 });
      },
      ...opts,
    });

    const second = new Transform({
      transform(chunk, encoding, callback) {
        callback(undefined, chunk['foo']);
      },
      ...opts,
    });

    const piped = pipeline(first, second);

    expect(piped.readableObjectMode).toStrictEqual(true);
    expect(piped.writableObjectMode).toStrictEqual(true);

    // Writing to first should have same effect as writing to piped (which should be duplex)
    first.write({ foo: 1 });
    piped.write({ foo: 2 });
    first.end();

    const result = await readAll(piped);

    expect(result).toEqual([3, 4]);
  });
});

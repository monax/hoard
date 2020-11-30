import { Readable, Transform, TransformOptions, Writable } from 'stream';
import { pipeline, readAll, waitFor } from './stream';

// This test is mostly to sanity check pumpify and to provide some executable documentation to remind you, dear reader,
// of its behaviour.
describe('streams', () => {
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
      transform(chunk, encoding, callback) {
        callback(undefined, chunk);
      },
      ...opts,
    });

    const second = new Transform({
      transform(chunk, encoding, callback) {
        callback(undefined, chunk);
      },
      ...opts,
    });

    const piped = pipeline(first, second);

    expect(piped.readableObjectMode).toStrictEqual(true);
    expect(piped.writableObjectMode).toStrictEqual(true);

    first.write({ foo: 1 });
    first.write({ foo: 2 });
    first.end();

    const result = await readAll(piped);

    expect(result).toEqual([{ foo: 1 }, { foo: 2 }]);
  });
});

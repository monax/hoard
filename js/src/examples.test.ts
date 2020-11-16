import {example} from "../src/examples";

describe('Examples should run', function () {
  test('buffer based plaintext', async function () {
    const isIIIIT = Buffer.from([1,2]) instanceof Uint8Array
    const data = Buffer.from('some stuff', 'utf8');
    const salt = Buffer.from('foo', 'ascii');
    await example(data, salt);
  })
  test('base64 based plaintext', async function () {
    const data = Buffer.from('some stuff', 'utf8').toString('base64');
    const salt = Buffer.from('foo', 'ascii').toString('base64');
    await example(data, salt);
  })
})

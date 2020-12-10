import { example } from '../src/examples';
import {openpgpExample} from "./openpgp";

describe('examples', function () {
  test('buffer based plaintext', async () => {
    const data = Buffer.from('some stuff', 'utf8');
    const salt = Buffer.from('foo', 'ascii');
    await example(data, salt);
  });

  test('base64 based plaintext', async () => {
    const data = Buffer.from('some stuff', 'utf8').toString('base64');
    const salt = Buffer.from('foo', 'ascii').toString('base64');
    await example(data, salt);
  });

  // TODO: hook this up for automated testing and update to work with array references
  // test('openpgp', async () => {
  //   const data = Buffer.from('some stuff', 'utf8').toString('base64');
  //   const salt = Buffer.from('foo', 'ascii').toString('base64');
  //   await openpgpExample(data, salt);
  // });
});

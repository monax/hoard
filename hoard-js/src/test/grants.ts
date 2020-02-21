import { 
  Client, 
  Read,
  ReadUntil,
  Write,
  Grant,
  NewHeader,
  NewPlaintext,
  NewPlaintextAndGrantSpec,
  ReducePlaintext,
  ReadLengthPrefixed,
  PlaintextToBytes,
  NewSymmetricSpec,
  PlaintextAndGrantSpec,
  ChunkData,
} from '../index';

import { Readable } from 'stream';
import assert = require('assert');
import * as fixtures from './fixtures';

describe('Should be able to store plaintext under symmetric grant', function () {
  const hoard = Client('localhost:53431')

  it('valid secret id', async function () {
    const data = Buffer.from('some stuff', 'utf8');
    const salt = Buffer.from('foo', 'ascii');
    const ptgs = [
      NewPlaintextAndGrantSpec(NewPlaintext(null, NewHeader(salt)), NewSymmetricSpec('testing-id')),
      NewPlaintextAndGrantSpec(NewPlaintext(data, null)),
    ];

    const grant = await Write<PlaintextAndGrantSpec, Grant>(ptgs, hoard.putSeal.bind(hoard));
    const decrypted = ReducePlaintext(await Read([], hoard.unsealGet(grant)));
    assert.strictEqual(Buffer.from(decrypted.getBody_asB64(), 'base64').toString(), data.toString())
    assert.strictEqual(Buffer.from(decrypted.getHead().getSalt_asB64(), 'base64').toString(), salt.toString())
  })

  it('rejects invalid secret id', async function () {
    const data = Buffer.from('some stuff', 'utf8');
    const ptgs = [
      NewPlaintextAndGrantSpec(NewPlaintext(), NewSymmetricSpec('bad-public-ID')),
      NewPlaintextAndGrantSpec(NewPlaintext(data)),
    ];

    // NOTE: we have no assert.rejects in node 9.7.1 which we currently need to support
    Write<PlaintextAndGrantSpec, Grant>(ptgs, hoard.putSeal.bind(hoard))
      .then(() => assert.fail('should fail when PublicID id not known to Hoard'))
      .catch(() => {}) // good
  })

  it('writes and reads stream', async () => {
    const veryLongText = Buffer.from(fixtures.LONG_TEXT.repeat(200), 'utf8');
    const chunks = ChunkData(veryLongText, 100)
    const ptgs = [
      NewPlaintextAndGrantSpec(NewPlaintext(null, NewHeader()), NewSymmetricSpec('testing-id')),
      ...chunks.map(chunk => NewPlaintextAndGrantSpec(NewPlaintext(chunk))),
    ];

    const grant = await Write<PlaintextAndGrantSpec, Grant>(ptgs, hoard.putSeal.bind(hoard));
    const plaintext = ReducePlaintext(await Read([], hoard.unsealGet(grant)));
    assert.strictEqual(Buffer.from(plaintext.getBody_asU8()).toString(), veryLongText.toString());
  }).timeout(5000)

  it('can stop reading a stream early', async () => {
    const veryLongText = Buffer.from(fixtures.LONG_TEXT.repeat(100), 'utf8');
    const salt = Buffer.from('foo', 'ascii');
    const chunks = ChunkData(veryLongText, 20)
    const firstN = 3;
    const ptgs = [
      NewPlaintextAndGrantSpec(NewPlaintext(null, NewHeader(salt)), NewSymmetricSpec('testing-id')),
      ...chunks.map(chunk => NewPlaintextAndGrantSpec(NewPlaintext(chunk))),
    ];

    const grant = await Write<PlaintextAndGrantSpec, Grant>(ptgs, hoard.putSeal.bind(hoard));
    const pts = await ReadUntil([], hoard.unsealGet(grant), (acc) => acc.length >= firstN ? true : false);
    assert.strictEqual(pts.length, firstN);
  })

  it('transforms bytes', async () => {
    const veryLongText = Buffer.from(fixtures.LONG_TEXT.repeat(100), 'utf8');
    const chunks = ChunkData(veryLongText, 100)
    const ptgs = [
      NewPlaintextAndGrantSpec(NewPlaintext(null, NewHeader()), NewSymmetricSpec('testing-id')),
      ...chunks.map(chunk => NewPlaintextAndGrantSpec(NewPlaintext(chunk))),
    ];

    const grant = await Write<PlaintextAndGrantSpec, Grant>(ptgs, hoard.putSeal.bind(hoard));
    const stream = hoard.unsealGet(grant);
    const ts = stream.pipe(PlaintextToBytes());

    const bytesToRead = 100;
    let bytes: Buffer = await new Promise((resolve, reject) => {
      stream.on("error", () => {});
      ts.on('readable', () => {
        let bs = ts.read(bytesToRead);
        stream.cancel();
        resolve(bs);
      });
    })

    assert.strictEqual(bytes.length, bytesToRead);
  })

  it('can read length prefixed value', async () => {
    // 5 is the single-byte length prefix, [1,2,3,4,5] should be value extracted
    const vals = [5, 1, 2, 3, 4, 5, 6, 5, 4, 3, 4, 5, 6, 7, 8, 8, 6, 4, 4];
    const stream = new Readable({ objectMode: true });
    vals.map(v => (NewPlaintext(Buffer.from([v]), null))).forEach(v => stream.push(v));
    stream.push(null);
    const plaintextStream = stream.pipe(PlaintextToBytes(2));
    const bs = await ReadLengthPrefixed(plaintextStream, 1);
    const expected = Buffer.from(vals.slice(1, vals[0] + 1));
    assert(bs.equals(expected), `readLengthPrefixed should read ${expected} but got ${bs}`)
  })
})

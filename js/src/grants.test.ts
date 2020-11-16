import {Readable} from 'stream';
import * as fixtures from './fixtures';
import {
  ChunkData,
  Client,
  NewHeader,
  NewPlaintext,
  NewPlaintextAndGrantSpec,
  NewSymmetricSpec,
  PlaintextToBytes,
  Read,
  ReadLengthPrefixed,
  ReducePlaintext,
} from './index';

describe('Grants', function () {
  const hoard = new Client('localhost:53431')

  test('valid secret id', async function () {
    const data = Buffer.from('some stuff', 'utf8');
    const salt = Buffer.from('foo', 'ascii');
    const ptgs = [
      NewPlaintextAndGrantSpec(NewPlaintext(undefined, NewHeader(salt)), NewSymmetricSpec('testing-id')),
      NewPlaintextAndGrantSpec(NewPlaintext(data, undefined)),
    ];

    const grant = await hoard.putSeal(ptgs);
    const decrypted = ReducePlaintext(await hoard.unsealGet(grant));
    expect(Buffer.from(decrypted.getBody_asB64(), 'base64').toString()).toStrictEqual(data.toString())
    const head = decrypted.getHead();
    if (!head) {
      throw new Error(`Should contain header`)
    }
    expect(Buffer.from(head.getSalt_asB64(), 'base64').toString()).toStrictEqual(salt.toString())
  })

  test('rejects invalid secret id', async function () {
    const data = Buffer.from('some stuff', 'utf8');
    const salt = Buffer.from('foo', 'ascii');
    const ptgs = [
      NewPlaintextAndGrantSpec(NewPlaintext(undefined, NewHeader(salt)), NewSymmetricSpec('bad-testing-id')),
      NewPlaintextAndGrantSpec(NewPlaintext(data)),
    ];

    // NOTE: we have no assert.rejects in node 9.7.1 which we currently need to support
    await expect(hoard.putSeal(ptgs)).rejects.toThrow()
  })

  test('writes and reads stream', async () => {
    const veryLongText = Buffer.from(fixtures.LONG_TEXT.repeat(200), 'utf8');
    const chunks = ChunkData(veryLongText, 100)
    const ptgs = [
      NewPlaintextAndGrantSpec(NewPlaintext(undefined, NewHeader()), NewSymmetricSpec('testing-id')),
      ...chunks.map(chunk => NewPlaintextAndGrantSpec(NewPlaintext(chunk))),
    ];

    const grant = await hoard.putSeal(ptgs);
    const plaintext = ReducePlaintext(await hoard.unsealGet(grant));
    expect(Buffer.from(plaintext.getBody_asU8()).toString()).toStrictEqual(veryLongText.toString());
  })

  test('can stop reading a stream early', async () => {
    const veryLongText = Buffer.from(fixtures.LONG_TEXT.repeat(100), 'utf8');
    const salt = Buffer.from('foo', 'ascii');
    const chunks = ChunkData(veryLongText, 20)
    const firstN = 3;
    const ptgs = [
      NewPlaintextAndGrantSpec(NewPlaintext(undefined, NewHeader(salt)), NewSymmetricSpec('testing-id')),
      ...chunks.map(chunk => NewPlaintextAndGrantSpec(NewPlaintext(chunk))),
    ];

    const grant = await hoard.putSeal(ptgs);
    const pts = await Read(hoard.grant.unsealGet(grant), (acc) => acc.length >= firstN);
    expect(pts.length).toStrictEqual(firstN);
  })

  test('transforms bytes', async () => {
    const veryLongText = Buffer.from(fixtures.LONG_TEXT.repeat(100), 'utf8');
    const chunks = ChunkData(veryLongText, 100)
    const ptgs = [
      NewPlaintextAndGrantSpec(NewPlaintext(undefined, NewHeader()), NewSymmetricSpec('testing-id')),
      ...chunks.map(chunk => NewPlaintextAndGrantSpec(NewPlaintext(chunk))),
    ];

    const grant = await hoard.putSeal(ptgs);
    const stream = hoard.grant.unsealGet(grant);
    const ts = stream.pipe(PlaintextToBytes());

    const bytesToRead = 100;
    const bytes: Buffer = await new Promise((resolve, reject) => {
      stream.on("error", () => null);
      ts.on('readable', () => {
        const bs = ts.read(bytesToRead);
        stream.cancel();
        resolve(bs);
      });
    })

    expect(bytes.length).toStrictEqual( bytesToRead);
  })

  test('can read length prefixed value', async () => {
    // 5 is the single-byte length prefix, [1,2,3,4,5] should be value extracted
    const vals = [5, 1, 2, 3, 4, 5, 6, 5, 4, 3, 4, 5, 6, 7, 8, 8, 6, 4, 4];
    const stream = new Readable({ objectMode: true });
    vals.map(v => (NewPlaintext(Buffer.from([v]), undefined))).forEach(v => stream.push(v));
    stream.push(null);
    const plaintextStream = stream.pipe(PlaintextToBytes(2));
    const bs = await ReadLengthPrefixed(plaintextStream, 1);
    const expected = Buffer.from(vals.slice(1, vals[0] + 1));
    expect(bs).toEqual(expected)
  })
})

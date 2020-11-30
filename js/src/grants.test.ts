import { Readable } from 'stream';
import * as fixtures from './fixtures';
import { Client, Header, make, Spec, SymmetricSpec } from './index';
import { bytesReadable, readAll, readBytes, readLengthPrefixed } from './stream';

describe('Grants', function () {
  const hoard = new Client('localhost:53431');
  const spec = make(Spec, (s) => s.setSymmetric(make(SymmetricSpec, (ss) => ss.setPublicid('testing-id'))));
  const badSpec = make(Spec, (s) => s.setSymmetric(make(SymmetricSpec, (ss) => ss.setPublicid('bad-testing-id'))));
  const emptyHeader = make(Header);

  test('valid secret id', async function () {
    const data = Buffer.from('some stuff', 'utf8');
    const salt = Buffer.from('foo', 'ascii');

    const header = make(Header, (h) => h.setSalt(salt));

    header.serializeBinary();
    const grant = await hoard.putSeal(spec, bytesReadable(data), header);
    const { head, body } = await hoard.unsealGet(grant);
    const decrypted = await readBytes(body);
    expect(decrypted.toString()).toStrictEqual(data.toString());
    if (!head) {
      throw new Error(`Should contain header`);
    }
    expect(Buffer.from(head.getSalt_asB64(), 'base64').toString()).toStrictEqual(salt.toString());
  });

  test('rejects invalid secret id', async function () {
    const data = Buffer.from('some stuff', 'utf8');
    const salt = Buffer.from('foo', 'ascii');

    await expect(hoard.putSeal(badSpec, data)).rejects.toThrow();
  });

  test('writes and reads stream', async () => {
    const veryLongText = fixtures.LONG_TEXT.repeat(1000);

    const grant = await hoard.putSeal(spec, veryLongText, emptyHeader);
    const { body } = await hoard.unsealGet(grant);
    const bs = await readBytes(body);
    expect(bs.toString()).toStrictEqual(veryLongText.toString());
  });

  test('can stop reading a stream early', async () => {
    const veryLongText = Buffer.from(fixtures.LONG_TEXT.repeat(100), 'utf8');
    const firstN = 3;

    const grant = await hoard.putSeal(spec, veryLongText);
    const pts = await readAll(hoard.grant.unsealGet(grant), (acc) => acc.length >= firstN);
    expect(pts.length).toStrictEqual(firstN);
  });

  test('can read length prefixed value', async () => {
    // 5 is the single-byte length prefix, [1,2,3,4,5] should be value extracted
    const vals = [5, 1, 2, 3, 4, 5, 6, 5, 4, 3, 4, 5, 6, 7, 8, 8, 6, 4, 4];
    const bs = await readLengthPrefixed(bytesReadable(vals), 1);
    const expected = Buffer.from(vals.slice(1, vals[0] + 1));
    expect(bs).toEqual(expected);
  });
});

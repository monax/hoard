import * as fixtures from './fixtures';
import { Client, deserializeGrant, Header, make, Plaintext, Spec, SymmetricSpec } from './index';
import { bytesReadable, readAll, readBytes } from './streaming';

const MiB = 1 << 20;

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
    const txt = fixtures.LONG_TEXT;
    const veryLongText = Buffer.from(txt.repeat((100 * MiB) / txt.length), 'utf8');
    const firstN = 2;

    const grant = await hoard.putSeal(spec, veryLongText);
    const pts = await readAll<Plaintext>(hoard.grant.unsealGet(grant), (acc) => acc.length >= firstN);
    expect(pts.length).toStrictEqual(firstN);
  });

  test('can put and get large text', async () => {
    const txt = fixtures.LONG_TEXT;
    const size = 200 * MiB;
    const veryLongText = Buffer.from(txt.repeat(size / txt.length), 'utf8');

    const grant = await hoard.putSeal(spec, veryLongText);
    const { body } = await hoard.unsealGet(grant);
    const output = await readBytes(body, veryLongText.length + 1);
    expect(output.length).toStrictEqual(veryLongText.length);
  });

  test('Stream produces error with an empty grant', async () => {
    // This gets pass the Protobuf binary decoder evenn though it is junk because the first byte (D == 12) is
    // interpreted as the 'End Group' wire type, because when masked with 7 (0b111), 7 & 12 = 4 = <End Group>
    const grant = deserializeGrant('DEADBEEFCAFEBABE');
    await expect(() => hoard.unsealGet(grant)).rejects.toThrow(/grant type not recognised/);
  });
});

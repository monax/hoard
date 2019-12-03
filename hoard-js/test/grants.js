const { Readable } = require('stream')
const Hoard = require('../index.js')
const assert = require('assert')
const fixtures = require('./fixtures')

describe('Should be able to store plaintext under symmetric grant', function () {
  const hoard = new Hoard.Client('localhost:53431')

  it('valid secret id', async function () {
    const text = 'some stuff'
    const salt = 'salty'
    const spec = { Symmetric: { PublicID: 'testing-id' } }

    const plaintextAndGrantSpec = Hoard.PlaintextAndGrantSpec(spec, text, salt)
    const grant = await Hoard.write(hoard.putSeal(), plaintextAndGrantSpec)
    const decrypted = await Hoard.read(hoard.unsealGet(grant), Hoard.Plaintext.reduce)
    assert.strictEqual(decrypted.Data.toString('utf8'), text)
    assert.strictEqual(decrypted.Salt.toString('utf8'), salt)
  })

  it('rejects invalid secret id', async function () {
    const text = 'some stuff'
    const spec = { Symmetric: { PublicID: 'bad-public-ID' } }
    const plaintextAndGrantSpec = Hoard.PlaintextAndGrantSpec(spec, text, 'foo')
    // NOTE: we have no assert.rejects in node 9.7.1 which we currently need to support
    Hoard.write(hoard.putSeal, plaintextAndGrantSpec)
      .then(() => assert.fail('should fail when PublicID id not known to Hoard'))
      .catch(() => {}) // good
  })

  it('writes and reads stream', async () => {
    const veryLongText = fixtures.LONG_TEXT.repeat(200)
    const spec = { Symmetric: { PublicID: 'testing-id' } }
    const grt = await Hoard.write(hoard.putSeal(), Hoard.PlaintextAndGrantSpec(spec, veryLongText).chunks(100))
    const pt = await Hoard.read(hoard.unsealGet(grt), Hoard.Plaintext.reduce)
    assert.strictEqual(pt.Data.toString(), veryLongText)
  }).timeout(5000)

  it('can stop reading a stream early', async () => {
    const veryLongText = fixtures.LONG_TEXT.repeat(1000)
    const spec = { Symmetric: { PublicID: 'testing-id' } }
    const grt = await Hoard.write(hoard.putSeal(), Hoard.PlaintextAndGrantSpec(spec, veryLongText, 'salt').chunks())

    const firstN = 14
    const pts = await Hoard.read(hoard.unsealGet(grt), (acc, val, returnNow) =>
      (acc.length < firstN) ? [...acc, val] : returnNow(acc)
    , [])
    assert.strictEqual(pts.length, firstN)
  })

  it('transforms bytes', async () => {
    const veryLongText = fixtures.LONG_TEXT.repeat(100)
    const spec = { Symmetric: { PublicID: 'testing-id' } }
    const grt = await Hoard.write(hoard.putSeal(), Hoard.PlaintextAndGrantSpec(spec, veryLongText).chunks(100))
    const stream = hoard.unsealGet(grt)
    const ts = stream.pipe(Hoard.ObjectsToBytes())

    const bytesToRead = 100
    let bs
    ts.on('readable', () => {
      bs = ts.read(bytesToRead)
      stream.close()
    })
    await stream.wait()
    assert.strictEqual(bs.length, bytesToRead)
  })

  it('can read length prefixed value', async () => {
    // 5 is the single-byte length prefix, [1,2,3,4,5] should be value extracted
    const vals = [5, 1, 2, 3, 4, 5, 6, 5, 4, 3, 4, 5, 6, 7, 8, 8, 6, 4, 4]
    const stream = new Readable({ objectMode: true })
    vals.map(v => ({ Data: Buffer.from([v]) })).forEach(v => stream.push(v))
    stream.push(null)
    const plaintextStream = stream.pipe(Hoard.ObjectsToBytes(msg => msg.Data,  2))
    const bs = await Hoard.readLengthPrefixed(plaintextStream, 1)
    const expected = Buffer.from(vals.slice(1, vals[0] + 1))
    assert(bs.equals(expected), `readLengthPrefixed should read ${expected} but got ${bs}`)
  })
})

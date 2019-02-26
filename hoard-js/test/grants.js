const Hoard = require('../index.js');
const assert = require('assert');

describe('Should be able to store plaintext under symmetric grant', function () {
    it('valid secret id', async function () {
        let text = 'some stuff'
        let salt = 'salty'
        let plaintextAndGrantSpec = {
            Plaintext: {
                Data: Buffer.from(text, 'utf8').toString('base64'),
                Salt: Buffer.from(salt, 'utf8').toString('base64')
            },
            GrantSpec: {
                Symmetric: {
                    PublicID: "test"
                }
            }
        }
        let hoard = new Hoard.Client('localhost:53431');
        let grant = await hoard.putseal(plaintextAndGrantSpec)
        let decrypted = await hoard.unsealget(grant)
        assert.strictEqual(decrypted.Data.toString('utf8'), text)
        assert.strictEqual(decrypted.Salt.toString('utf8'), salt)
    });

    it('invalid secret id', async function () {
        let text = 'some stuff'
        let plaintextAndGrantSpec = {
            Plaintext: {
                Data: Buffer.from(text, 'utf8').toString('base64'),
                Salt: Buffer.from('foo', 'ascii').toString('base64')
            },
            GrantSpec: {
                Symmetric: {
                    PublicID: "bad-public-ID"
                }
            }
        }
        let hoard = new Hoard.Client('localhost:53431');
        assert.rejects(() => hoard.putseal(plaintextAndGrantSpec), Error,
            "should fail when PublicID id not known to Hoard");
    });
});

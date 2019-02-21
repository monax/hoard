let Hoard = require('./index.js');

var fs = require('fs');
const path = require('path')

// npm install openpgp
const openpgp = require('openpgp');
openpgp.initWorker({ path:'openpgp.worker.js' })

let hoard = new Hoard.Client('localhost:53431');
let plaintext = {
    Data: Buffer.from('some stuff', 'utf8'),
    Salt: Buffer.from('foo', 'ascii')
};

// Below is an example of using Hoard's asymmetric grants
// Note: This will only work if the hoard daemon is configured to do PGP signing
const example = async function (plaintextIn) {
    try {
        var plaintext, ref, grant

        var pubkey = fs.readFileSync(path.join(__dirname, '../grant', 'public.key.asc'));
        var privkey = fs.readFileSync(path.join(__dirname, '../grant', 'private.key.asc'));
    
        grantIn = {
            Plaintext: plaintextIn,
            GrantSpec: {
                OpenPGP: {
                    PublicKey: pubkey.toString()
                }
            }
        }
    
        grant = await hoard.putseal(grantIn);
        console.log(grant.Data)
        console.log(hoard.base64ify(grant));

        const options = {
            message: await openpgp.message.readArmored(grant.EncryptedReference),
            privateKeys: [(await openpgp.key.readArmored(privkey)).keys[0]]
        }

        ref = JSON.parse((await openpgp.decrypt(options)).data)
        console.log(hoard.base64ify(ref));
        plaintext = await hoard.get(ref)
        console.log(plaintext)
        console.log('Plaintext (PGP Grant): ' + plaintext.Data.toString());
    }
    catch (err) {
        console.log(err)
    }
};

example(plaintext);

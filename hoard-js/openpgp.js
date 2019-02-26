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

        // For this example we will lock under our own keypair
        var pubkey = fs.readFileSync(path.join(__dirname, '../grant', 'public.key.asc'));
        var privkey = fs.readFileSync(path.join(__dirname, '../grant', 'private.key.asc'));
    
        // The openpgp grant allows us to encrypt the reference
        // under the specified public key
        grantIn = {
            Plaintext: plaintextIn,
            GrantSpec: {
                OpenPGP: {
                    PublicKey: pubkey.toString()
                }
            }
        }
    
        grant = await hoard.putseal(grantIn);
        console.log(hoard.base64ify(grant));

        // A locally running Hoard daemon can also be configured to decrypt
        // the given grant based on a local keyring.
        // For this example, we will use Node's OpenPGP library.
        const options = {
            message: await openpgp.message.readArmored(grant.EncryptedReference),
            privateKeys: [(await openpgp.key.readArmored(privkey)).keys[0]]
        };

        ref = JSON.parse((await openpgp.decrypt(options)).data);
        console.log(hoard.base64ify(ref));
        plaintext = await hoard.get(ref);
        console.log(plaintext);
        console.log('Plaintext (PGP Grant): ' + plaintext.Data.toString());
    }
    catch (err) {
        console.log(err);
        process.exit(1);
    }
};

example(plaintext);

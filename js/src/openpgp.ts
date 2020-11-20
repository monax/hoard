import { 
    Client,
    Write,
    Duplex,
    Grant,
    ReducePlaintext,
    NewHeader,
    NewPlaintext,
    NewPlaintextAndGrantSpec,
    NewOpenPGPSpec,
    PlaintextAndGrantSpec,
} from './index';

import * as fs from 'fs';
import * as path from 'path';

const hoard = Client('localhost:53431');

// yarn install openpgp
const openpgp = require('openpgp');
openpgp.initWorker({ path:'openpgp.worker.js' })

const data = Buffer.from('some stuff', 'utf8');
const salt = Buffer.from('foo', 'ascii');

// Below is an example of using Hoard's asymmetric grants
// Note: This will only work if the hoard daemon is configured to do PGP signing
const example = async function (data: string | Uint8Array, salt: string | Uint8Array) {
    try {
        let plaintexts = [NewPlaintext(data, NewHeader(salt))];
    
        // For this example we will lock under our own keypair
        var pubkey = fs.readFileSync(path.join(__dirname, '../grant', 'public.key.asc'));
        var privkey = fs.readFileSync(path.join(__dirname, '../grant', 'private.key.asc'));
    
        // The openpgp grant allows us to encrypt the reference
        // under the specified public key
        const ptgs = plaintexts.map(pt => NewPlaintextAndGrantSpec(pt, NewOpenPGPSpec(pubkey.toString())));
        const grant = await Write<PlaintextAndGrantSpec, Grant>(ptgs, hoard.putSeal.bind(hoard))
        console.log(grant.toObject());

        // A locally running Hoard daemon can also be configured to decrypt
        // the given grant based on a local keyring.
        // For this example, we will use Node's OpenPGP library.
        const options = {
            message: await openpgp.message.readArmored(grant.getEncryptedreferences()),
            privateKeys: [(await openpgp.key.readArmored(privkey)).keys[0]]
        };

        let references = JSON.parse((await openpgp.decrypt(options)).data);
        console.log(references);
        plaintexts = await Duplex(references, [], hoard.get())
        console.log(ReducePlaintext(plaintexts));
    }
    catch (err) {
        console.log(err);
        process.exit(1);
    }
};

example(data, salt);

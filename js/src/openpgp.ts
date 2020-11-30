import * as fs from 'fs';
import { decrypt, DecryptOptions, initWorker, key, message } from 'openpgp';
import * as path from 'path';
import { Client, Header, make, OpenPGPSpec, Spec } from './index';
import { readBytes } from './stream';

const hoard = new Client('localhost:53431');

initWorker({ path: 'openpgp.worker.js' });

// Below is an example of using Hoard's asymmetric grants
// Note: This will only work if the hoard daemon is configured to do PGP signing
export async function openpgpExample(data: string | Uint8Array, salt: string | Uint8Array): Promise<void> {
  // For this example we will lock under our own keypair
  const grantPath = '../../grant';
  const pubkey = fs.readFileSync(path.join(__dirname, grantPath, 'public.key.asc'));
  const privkey = fs.readFileSync(path.join(__dirname, grantPath, 'private.key.asc'));

  const head = make(Header);
  const spec = make(Spec, (s) => s.setOpenpgp(make(OpenPGPSpec, (ps) => ps.setPublickey(pubkey.toString()))));

  // The openpgp grant allows us to encrypt the reference
  // under the specified public key
  const grant = await hoard.putSeal(spec, data);

  // A locally running Hoard daemon can also be configured to decrypt
  // the given grant based on a local keyring.
  // For this example, we will use Node's OpenPGP library.
  const refs = JSON.parse(Buffer.from(grant.getEncryptedreferences_asU8()).toString('utf8'));
  const secretKey = refs['Refs'][0]['SecretKey'];
  console.log(secretKey)
  const ref = Buffer.from(secretKey, 'base64').toString('utf8');
  console.log(ref)
  const options: DecryptOptions = {
    message: await message.readArmored(refs),
    privateKeys: [(await key.readArmored(privkey)).keys[0]],
  };

  const references = JSON.parse((await decrypt(options)).data);
  console.log(references);
  const { body } = await hoard.get(references);
  const plaintext = await readBytes(body);
  console.log(plaintext);
}

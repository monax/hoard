// This gives a walk-through of using the Hoard API. All index.js of this
// library does it wraps the dynamically generated GRPC client in promises
// and and abstracts away the loading of the protobuf file. You may prefer
// to copy in the code and the hoard.proto file in order to communicate with
// the Hoard daemon.
import * as assert from 'assert';
import { Address, Client, Header, make, PlaintextSpec, Spec, SymmetricSpec } from './index';
import { readAll, readBytes } from './streaming';

// This is the default tcp socket that hoard runs on, just run `bin/hoard`
// after running `make build` in the main hoard repo.
const hoard = new Client('localhost:53431');

// All input and outputs to the API methods are JSON objects representing the
// message type with the parameters contained within. This corresponds to the
// `message` declarations in hoard.proto which can be used as reference.

// Below is an example of running through a series of hoard calls wrapped in promises.
// By wrapping this in an async function we can use await/async try/catch syntactic sugar around
export async function example(data: string | Uint8Array, salt: string | Uint8Array): Promise<void> {
  const head = make(Header, (h) => h.setSalt(salt));

  const assertBytesEqual = (actual: string | Uint8Array | undefined, expected: string | Uint8Array | undefined) => {
    if (!actual) {
      throw new Error(`actual is void`);
    }
    if (!expected) {
      throw new Error(`expected is void`);
    }
    actual = Buffer.from(actual);
    expected = Buffer.from(expected);
    if (Buffer.compare(actual, expected) !== 0) {
      throw new Error(`Expected body '${expected}' but actually got '${actual}'`);
    }
  };

  // Both the address and secret key are a deterministic function of the
  // data and the salt (the plaintext). You need the salt and secret key
  // to decrypt (or get).

  // console.log(plaintext.toObject())
  // Put the plaintext in storage
  let references = await readAll(hoard.put(data, head));
  assert.strictEqual(references.length, 2);

  // We can get the plaintext back by `get`ing the grant
  let actual = await hoard.get(references);
  assert.deepStrictEqual(actual?.head?.toObject(), head.toObject());
  assertBytesEqual(await readBytes(actual.body), data);

  // This time we'll just encrypt and ask for the result rather than storing it
  // We get a 'hypothetical' reference (since it is not stored) and he ciphertext itself
  const refAndCiphertexts = await readAll(hoard.encrypt(data, head));
  assert.deepStrictEqual(refAndCiphertexts[0].getReference()?.toObject(), references[0].toObject());

  // decrypt is our inverse
  // We can also use the ReadAll helper that by default will use the first object as accumulator
  actual = await hoard.decrypt(refAndCiphertexts);
  assertBytesEqual(await readBytes(actual.body), data);

  // Put it back to get a reference
  references = await readAll(hoard.put(data, head));

  // We can ask for file information (we could have just provided the grant here, but address is all that is needed)
  const statInfo = await hoard.stat(references[0]);
  assert.strictEqual(statInfo.getExists(), true);

  // Note that all arguments take an object, representing the message, so 'address' is {address: address}
  // pull interacts with underlying storage directly so fetches ciphertext
  const ciphertexts = await readBytes(hoard.pull([make(Address, (a) => a.setAddress(statInfo.getAddress_asU8()))]));
  await assertBytesEqual(ciphertexts, refAndCiphertexts[0].getCiphertext()?.getEncrypteddata_asU8());

  const addresses = await readAll(hoard.push(ciphertexts));
  assert.deepStrictEqual(addresses[0].getAddress_asU8(), references[0].getAddress_asU8());

  // A plaintext grant allows us to reference the reference without
  // encryption for ease of later retrieval

  let grant = await hoard.putSeal(
    make(Spec, (s) => s.setPlaintext(make(PlaintextSpec))),
    data,
    head,
  );
  assert.ok(grant);

  // We can get the plaintext back by `unsealGet`ing the grant
  actual = await hoard.unsealGet(grant);
  assertBytesEqual(await readBytes(actual.body), data);

  assert.strictEqual(actual.head?.getSalt_asB64(), 'Zm9v');

  // A symmetric grant allows us to encrypt the reference
  // through secrets configured on the hoard daemon
  grant = await hoard.putSeal(
    make(Spec, (s) => s.setSymmetric(make(SymmetricSpec, (ss) => ss.setPublicid('testing-id')))),
    data,
    head,
  );
  assert.ok(grant);

  actual = await hoard.unsealGet(grant);
  await assertBytesEqual(await readBytes(actual.body), data);

  // DELETE!
  const [deleted] = await readAll(hoard.unsealDelete(grant));

  // We will have deleted the _link_ so get that now
  const [linkRef] = await readAll(hoard.unseal(grant))

  assert.strictEqual(deleted.getAddress_asB64(), linkRef.getAddress_asB64());
}

// To run the async example in this case ignoring the promise result uncomment the statements below

// Lets store some data. Here we use a salt that means that we will get
// different bytes for our encryption that is semantically secure in the
// length of the salt. This is useful if we want to disguise that a known
// piece of text has been stored since it will give it a different address

// const data = Buffer.from('some stuff');
// const salt = Buffer.from('foo');
// example(data, salt);

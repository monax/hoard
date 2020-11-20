// This gives a walk-through of using the Hoard API. All index.js of this
// library does it wraps the dynamically generated GRPC client in promises
// and and abstracts away the loading of the protobuf file. You may prefer
// to copy in the code and the hoard.proto file in order to communicate with
// the Hoard daemon.
import {
  Client,
  NewPlaintext,
  NewAddress,
  ReducePlaintext,
  ReduceReferenceAndCiphertext,
  ReduceCiphertext,
  Plaintext,
  NewPlaintextAndGrantSpec,
  StatInfo,
  Grant,
  NewPlaintextSpec,
  NewSymmetricSpec,
  PlaintextAndGrantSpec,
  Header,
  ReadHeader,
} from "./index";
import * as assert from "assert";

// This is the default tcp socket that hoard runs on, just run `bin/hoard`
// after running `make build` in the main hoard repo.
const hoard = new Client("localhost:53431");

// All input and outputs to the API methods are JSON objects representing the
// message type with the parameters contained within. This corresponds to the
// `message` declarations in hoard.proto which can be used as reference.

// Below is an example of running through a series of hoard calls wrapped in promises.
// By wrapping this in an async function we can use await/async try/catch syntactic sugar around
export async function example(
  data: string | Uint8Array,
  salt: string | Uint8Array
): Promise<void> {
  const head = new Header();
  head.setSalt(salt);
  const expected = NewPlaintext(data, head);
  let plaintext = expected;

  // Both the address and secret key are a deterministic function of the
  // data and the salt (the plaintext). You need the salt and secret key
  // to decrypt (or get).

  // console.log(plaintext.toObject())
  // Put the plaintext in storage
  let references = await hoard.put(plaintext);
  assert.strictEqual(references.length, 2);

  // We can get the plaintext back by `get`ing the grant
  let plaintexts = await hoard.get(references);
  plaintext = ReducePlaintext(plaintexts);
  assert.deepStrictEqual(plaintext.toObject(), expected.toObject());

  // This time we'll just encrypt and ask for the result rather than storing it
  // We get a 'hypothetical' reference (since it is not stored) and the ciphertext itself
  const refAndCiphertexts = await hoard.encrypt(plaintext);
  const refAndCiphertext = ReduceReferenceAndCiphertext(refAndCiphertexts);
  assert.deepStrictEqual(
    refAndCiphertext.getReference()?.toObject(),
    references[0].toObject()
  );

  // decrypt is our inverse
  // We can also use the ReadAll helper that by default will use the first object as accumulator
  plaintexts = await hoard.decrypt(refAndCiphertexts);
  plaintext = ReducePlaintext(plaintexts);
  assert.deepStrictEqual(plaintext.toObject(), expected.toObject());

  // Put it back to get a reference
  references = await hoard.put(plaintext);

  // We can ask for file information (we could have just provided the grant here, but address is all that is needed)
  const statInfo = await hoard.stat(references[0]);

  assert.strictEqual(statInfo.getExists(), true);

  // Note that all arguments take an object, representing the message, so 'address' is {address: address}
  // pull interacts with underlying storage directly so fetches ciphertext
  const ciphertexts = await hoard.pull([
    NewAddress(statInfo.getAddress_asU8()),
  ]);
  assert.deepStrictEqual(
    ReduceCiphertext(ciphertexts).getEncrypteddata_asU8(),
    refAndCiphertexts[0].getCiphertext()?.getEncrypteddata_asU8()
  );

  const addresses = await hoard.push(ciphertexts);
  assert.deepStrictEqual(
    addresses[0].getAddress_asU8(),
    references[0].getAddress_asU8()
  );

  // A plaintext grant allows us to reference the reference without
  // encryption for ease of later retrieval

  let ptgs = [
    NewPlaintextAndGrantSpec(NewPlaintext(undefined, head), NewPlaintextSpec()),
    NewPlaintextAndGrantSpec(NewPlaintext(data)),
  ];
  let grant = await hoard.putSeal(ptgs);
  assert.ok(grant);

  // We can get the plaintext back by `unsealget`ing the grant
  plaintexts = await hoard.unsealGet(grant);
  plaintext = ReducePlaintext(plaintexts);
  assert.deepStrictEqual(plaintext.toObject(), expected.toObject());

  const header = await ReadHeader(hoard.grant.unsealGet(grant));
  assert.strictEqual(header.getSalt_asB64(), "Zm9v");

  // A symmetric grant allows us to encrypt the reference
  // through secrets configured on the hoard daemon
  ptgs = [
    NewPlaintextAndGrantSpec(
      NewPlaintext(undefined, head),
      NewSymmetricSpec("testing-id")
    ),
    NewPlaintextAndGrantSpec(NewPlaintext(data)),
  ];
  grant = await hoard.putSeal(ptgs);
  assert.ok(grant);

  plaintexts = await hoard.unsealGet(grant);
  plaintext = ReducePlaintext(plaintexts);
  assert.deepStrictEqual(plaintext.toObject(), expected.toObject());

  const deleted = await hoard.unsealDelete(grant);
  assert.deepStrictEqual(deleted[0], addresses[0]);
}

// To run the async example in this case ignoring the promise result uncomment the statements below

// Lets store some data. Here we use a salt that means that we will get
// different bytes for our encryption that is semantically secure in the
// length of the salt. This is useful if we want to disguise that a known
// piece of text has been stored since it will give it a different address

// const data = Buffer.from('some stuff');
// const salt = Buffer.from('foo');
// example(data, salt);

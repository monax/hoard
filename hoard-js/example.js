// This gives a walk-through of using the Hoard API. All index.js of this
// library does it wraps the dynamically generated GRPC client in promises
// and and abstracts away the loading of the protobuf file. You may prefer
// to copy in the code and the hoard.proto file in order to communicate with
// the Hoard daemon.
let Hoard = require('./index.js');
var messages = require('./protobuf/hoard_pb');

// This is the default tcp socket that hoard runs on, just run `bin/hoard`
// after running `make build` in the main hoard repo.
let hoard = new Hoard.Client('localhost:53431');

// All input and outputs to the API methods are JSON objects representing the
// message type with the parameters contained within. This corresponds to the
// `message` declarations in hoard.proto which can be used as reference.

// Lets's store some data. Here we use a salt that means that we will get
// different bytes for our encryption that is semantically secure in the
// length of the salt. This is useful if we want to disguise that a known
// piece of text has been stored since it will give it a different address
let plaintext = {
    Data: Buffer.from('some stuff', 'utf8'),
    Salt: Buffer.from('foo', 'ascii')
};

// Below is an example of running through a series of hoard calls wrapped in promises.
// By wrapping this in an async function we can use await/async try/catch syntactic sugar around
//
const example = async function (plaintextIn) {
    try {
        // Both the address and secret key are a deterministic function of the
        // data and the salt (the plaintext). You need the salt and secret key
        // to decrypt (or get).
        // Put the plaintext in storage
        let ref = await hoard.put(plaintextIn);
        // (Base64 should be standard text representation for address, secretKey, and salt)
        console.log(base64ify(ref));
        // We can get the plaintext back by `get`ing the reference
        let plaintext = await hoard.get(ref);
        console.log('Plaintext: ' + plaintext.Data.toString());
        // This time we'll just encrypt and ask for the result rather than storing it
        let refAndCiphertext = await hoard.encrypt(plaintext);
        // We get a 'hypothetical' reference (since it is not stored) and the ciphertext itself
        console.log(refAndCiphertext);
        // decrypt is our inverse
        let plaintext2 = await hoard.decrypt(refAndCiphertext);
        console.log('Plaintext (again): ' + plaintext2.Data.toString());
        // Put it back to get a ref
        let ref2 = await hoard.put(plaintext2);
        // We can ask for file information (we could have just provided the ref here, but address is all that is needed)
        let statInfo = await hoard.stat({Address: ref2.Address});
        console.log(statInfo);
        // Note that all arguments take an object, representing the message, so 'address' is {address: address}
        // pull interacts with underlying storage directly so fetches ciphertext
        let ciphertext2 = await hoard.pull({Address: statInfo.Address});
        console.log(ciphertext2);
        let address = await hoard.push(ciphertext2);
        console.log(address)
    }
    catch (err) {
        console.log(err)
    }
};

// Run the async example in this case ignoring the promise result
example(plaintext);

// Utility for printing message types
const base64ify = function (obj) {
    let newObj = {};
    for (let key of Object.keys(obj)) {
        let value = obj[key];
        if (value instanceof Buffer) {
            newObj[key] = value.toString(`base64`)
        } else {
            newObj[key] = value
        }
    }
    return newObj;
};
// This gives a walk-through of using the Hoard API. All index.js of this
// library does it wraps the dynamically generated GRPC client in promises
// and and abstracts away the loading of the protobuf file. You may prefer
// to copy in the code and the hoard.proto file in order to communicate with
// the Hoard daemon.
let Hoard = require('./index.js');

// This is the default tcp socket that hoard runs on, just run `bin/hoard`
// after running `make build` in the main hoard repo.
let hoard = new Hoard.Client('localhost:53431');

// All input and outputs to the API methods are JSON objects representing the
// message type with the parameters contained within. This corresponds to the
// `message` declarations in hoard.proto which can be used as reference.

// Below is an example of running through a series of hoard calls wrapped in promises.
// By wrapping this in an async function we can use await/async try/catch syntactic sugar around
const example = async function (plaintextIn) {
    try {
        var plaintext, ref, grant;

        // Both the address and secret key are a deterministic function of the
        // data and the salt (the plaintext). You need the salt and secret key
        // to decrypt (or get).
        // Put the plaintext in storage
        ref = await hoard.put(plaintextIn);
        // (Base64 should be standard text representation for address, secretKey, and salt)
        console.log(hoard.base64ify(ref));
        // We can get the plaintext back by `get`ing the grant
        plaintext = await hoard.get(ref);
        console.log('Plaintext (Reference): ' + plaintext.Data.toString());
        
        // This time we'll just encrypt and ask for the result rather than storing it
        let refAndCiphertext = await hoard.encrypt(plaintext);
        // We get a 'hypothetical' reference (since it is not stored) and the ciphertext itself
        console.log(refAndCiphertext);
        // decrypt is our inverse
        plaintext = await hoard.decrypt(refAndCiphertext);
        console.log('Plaintext (Decrypted): ' + plaintext.Data.toString());

        // Put it back to get a reference
        ref = await hoard.put(plaintext);
        // We can ask for file information (we could have just provided the grant here, but address is all that is needed)
        let statInfo = await hoard.stat({Address: ref.Address});
        console.log(statInfo);
        // Note that all arguments take an object, representing the message, so 'address' is {address: address}
        // pull interacts with underlying storage directly so fetches ciphertext
        let ciphertext = await hoard.pull({Address: statInfo.Address});
        console.log(ciphertext);
        let address = await hoard.push(ciphertext);
        console.log(address);

        // A plaintext grant allows us to reference the reference without
        // encryption for ease of later retrieval
        let grantIn = {
            Plaintext: plaintextIn,
            GrantSpec: {
                Plaintext: {}
            }
        };

        grant = await hoard.putseal(grantIn);
        console.log(hoard.base64ify(grant));

        // We can get the plaintext back by `unsealget`ing the grant
        plaintext = await hoard.unsealget(grant);
        console.log('Plaintext (Grant): ' + plaintext.Data.toString());

        // A symmetric grant allows us to encrypt the reference
        // through secrets configured on the hoard daemon
        grantIn = {
            Plaintext: plaintextIn,
            GrantSpec: {
                Symmetric: {
                    PublicID: Buffer.from('test', 'utf8')
                }
            }
        };

        grant = await hoard.putseal(grantIn);

        // Convert to string and back again
        grant = JSON.stringify(hoard.base64ify(grant));
        console.log(grant);
        grant = JSON.parse(grant);
        console.log(grant);

        plaintext = await hoard.unsealget(grant);
        console.log('Plaintext (Grant): ' + plaintext.Data.toString());
    }
    catch (err) {
        console.log(err);
        process.exit(1);
    }
};

exports.example = example

// To run the async example in this case ignoring the promise result uncomment the statements below

// Lets store some data. Here we use a salt that means that we will get
// different bytes for our encryption that is semantically secure in the
// length of the salt. This is useful if we want to disguise that a known
// piece of text has been stored since it will give it a different address

// let plaintext = {
//     Data: Buffer.from('some stuff', 'utf8'),
//     Salt: Buffer.from('foo', 'ascii')
// };

// example(plaintext);

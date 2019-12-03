// This gives a walk-through of using the Hoard API. All index.js of this
// library does it wraps the dynamically generated GRPC client in promises
// and and abstracts away the loading of the protobuf file. You may prefer
// to copy in the code and the hoard.proto file in order to communicate with
// the Hoard daemon.
const { Client, base64ify, read, write } = require('./index.js')

// This is the default tcp socket that hoard runs on, just run `bin/hoard`
// after running `make build` in the main hoard repo.
const hoard = new Client('localhost:53431')

// All input and outputs to the API methods are JSON objects representing the
// message type with the parameters contained within. This corresponds to the
// `message` declarations in hoard.proto which can be used as reference.

// Below is an example of running through a series of hoard calls wrapped in promises.
// By wrapping this in an async function we can use await/async try/catch syntactic sugar around
const example = async function (plaintextIn) {
  try {
    let plaintext = {}
    let ref
    let grant

    // Both the address and secret key are a deterministic function of the
    // data and the salt (the plaintext). You need the salt and secret key
    // to decrypt (or get).
    // Put the plaintext in storage
    let stream = await hoard.put()
    stream.write(plaintextIn)
    ref = await stream.close()
    // (Base64 should be standard text representation for address, secretKey, and salt)
    console.log(base64ify(ref))
    // We can get the plaintext back by `get`ing the grant
    stream = hoard.get(ref)
    stream.on('data', pt => Object.assign(plaintext, pt))
    await stream.wait()
    console.log('Plaintext (Reference): ' + plaintext.Data.toString())

    // This time we'll just encrypt and ask for the result rather than storing it
    stream = hoard.encrypt()
    stream.write(plaintext)
    const refAndCiphertext = await stream.close()
    // We get a 'hypothetical' reference (since it is not stored) and the ciphertext itself
    console.log(base64ify(refAndCiphertext))
    // decrypt is our inverse
    // We can also use the readAll helper that by default will use the first object as accumulator
    plaintext = await read(hoard.decrypt(refAndCiphertext), Object.assign)
    console.log('Plaintext (Decrypted): ' + plaintext.Data.toString())

    // Put it back to get a reference
    // we can use writeAll to write a number of component objects (or just one) to a stream then close
    ref = await write(hoard.put(), plaintext)
    // We can ask for file information (we could have just provided the grant here, but address is all that is needed)
    const statInfo = await hoard.stat({ Address: ref.Address })
    console.log(statInfo)
    // Note that all arguments take an object, representing the message, so 'address' is {address: address}
    // pull interacts with underlying storage directly so fetches ciphertext
    const ciphertext = await read(hoard.pull({ Address: statInfo.Address }), Object.assign)
    console.log(ciphertext)
    const address = await write(hoard.push(), ciphertext)
    console.log(address)

    // A plaintext grant allows us to reference the reference without
    // encryption for ease of later retrieval
    let plaintextAndGrantSpec = {
      Plaintext: plaintextIn,
      GrantSpec: {
        Plaintext: {}
      }
    }

    grant = await write(hoard.putSeal(), plaintextAndGrantSpec)
    console.log(base64ify(grant))

    // We can get the plaintext back by `unsealget`ing the grant
    plaintext = await read(hoard.unsealGet(grant), Object.assign)
    console.log('Plaintext (Grant): ' + plaintext.Data.toString())

    // A symmetric grant allows us to encrypt the reference
    // through secrets configured on the hoard daemon
    plaintextAndGrantSpec = {
      Plaintext: plaintextIn,
      GrantSpec: {
        Symmetric: {
          PublicID: Buffer.from('testing-id', 'utf8')
        }
      }
    }

    grant = await write(hoard.putSeal(), plaintextAndGrantSpec)

    // Convert to string and back again
    grant = JSON.stringify(base64ify(grant))
    console.log(grant)
    grant = JSON.parse(grant)
    console.log(grant)

    plaintext = await read(hoard.unsealGet(grant), Object.assign)
    console.log('Plaintext (Grant): ' + plaintext.Data.toString())

    const addr = await hoard.unsealDelete(grant)
    console.log('Deleted address: ' + addr.Address.toString('hex'))

    const plaintextAndGrantSpecAndMeta = {
      Meta: {
        Name: 'test'
      },
      PlaintextAndGrantSpec: {
        Plaintext: plaintextIn,
        GrantSpec: {
          Plaintext: {}
        }
      }
    }

    grant = await write(hoard.upload(), plaintextAndGrantSpecAndMeta)
    console.log(grant)

    const plaintextAndMeta = await read(hoard.download(grant), Object.assign)
    console.log(plaintextAndMeta)
  } catch (err) {
    console.log(err)
    process.exit(1)
  }
}

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

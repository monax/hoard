'use strict'

const path = require('path')
const { Transform } = require('stream')

const LOCAL_PATH = path.join(__dirname, './protobuf')
const PROTO_PATH = path.join(__dirname, '../protobuf')
const PROTO_FILE = 'api.proto'

const DEFAULT_CHUNK_SIZE = 2 ** 16

const protoLoader = require('@grpc/proto-loader')
const grpc = require('grpc')

const options = {
  keepCase: true,
  longs: String,
  enums: String,
  defaults: true,
  oneofs: true,
  includeDirs: [
    LOCAL_PATH, PROTO_PATH
  ]
}

function isIterable (obj) {
  return typeof Object(obj)[Symbol.iterator] === 'function'
}

const packageDefinition = protoLoader.loadSync(PROTO_FILE, options)

function Client (hoardURL) {
  if (!hoardURL) {
    throw new Error('Hoard Client requires the hoardURL be passed to the constructor')
  }
  const self = this
  const api = grpc.loadPackageDefinition(packageDefinition).api

  function lowerCamel (str) {
    return str.charAt(0).toLowerCase() + str.slice(1)
  }

  // Promisifies GRPC methods - should work with unary calls and client-side, server-side, and bidirectional streams
  function promisify (serviceClient, methodName) {
    const method = serviceClient[methodName].bind(serviceClient)
    const { requestStream, responseStream } = serviceClient.$method_definitions[methodName]

    return (...args) => {
      // Capture streamOrValue from inside promise
      let streamOrValue

      const promise = new Promise((resolve, reject) => {
        // Call the GRPC method
        streamOrValue = method(...[...args, (error, value) => error ? reject(error) : resolve(value)])
        streamOrValue.on('error', err => err.code === grpc.status.CANCELLED ? resolve(err) : reject(err))
        streamOrValue.on('close', () => resolve())
        streamOrValue.on('end', () => resolve())
      })

      if (!requestStream && !responseStream) {
        // method is unary
        // Caller just needs to wait until their callback will not longer be called
        return promise
      }

      // requestStream === true => method is client-side or bi-directional streaming...
      // Override end function to return the result promise after calling original end
      streamOrValue.close = (...endArgs) => {
        // Server-side stream
        if (requestStream) {
          streamOrValue.end(...endArgs)
        }
        // client-side stream
        if (responseStream) {
          streamOrValue.cancel()
        }
        return promise
      }

      streamOrValue.wait = () => promise
      // Maintain consistent promise return
      return streamOrValue
    }
  }

  // GRPC service definitions to include on Hoard client
  const ServiceClients = [
    api.Cleartext,
    api.Encryption,
    api.Storage,
    api.Grant,
    api.Document
  ]
  // Smoosh all the methods from different services together
  for (const ServiceClient of ServiceClients) {
    const client = new ServiceClient(hoardURL, grpc.credentials.createInsecure())

    for (const methodName of Object.keys(client.$method_definitions)) {
      // Create an async version of each function and add it to the mono client
      self[lowerCamel(methodName)] = promisify(client, methodName)
    }
  }
}

function Plaintext (data, salt) {
  const newBuffer = d => d ? Buffer.from(d) : Buffer.alloc(0)
  const msg = {
    Data: newBuffer(data),
    Salt: newBuffer(salt)
  }

  // Generate chunks of Plaintext suitable for streaming
  msg.chunks = function * (chunkSize = DEFAULT_CHUNK_SIZE) {
    let data = Buffer.from(msg.Data)
    yield { Salt: msg.Salt }
    while (data.length > 0) {
      yield { Data: data.slice(0, chunkSize) }
      data = data.slice(chunkSize)
    }
  }

  return msg
}

Plaintext.reduce = function (acc, val) {
  return {
    Salt: acc.Salt.length ? acc.Salt : val.Salt,
    Data: Buffer.concat([acc.Data, val.Data])
  }
}

function PlaintextAndGrantSpec (spec, data, salt) {
  const msg = {
    Plaintext: Plaintext(data, salt),
    GrantSpec: spec
  }

  msg.chunks = function * (chunkSize = DEFAULT_CHUNK_SIZE) {
    yield { GrantSpec: msg.GrantSpec }
    for (const pt of msg.Plaintext.chunks(chunkSize)) {
      yield { Plaintext: pt }
    }
  }

  return msg
}

PlaintextAndGrantSpec.reduce = function (acc, val) {
  return {
    Plaintext: acc.Plaintext.merge(val.Plaintext),
    GrantSpec: acc.GrantSpec || val.GrantSpec
  }
}

// Converts an object mode stream to a byte mode stream selecting a single buffer from the source using dataSelector
// and buffering bytes in chunks of at chunkSize
const BufferTransform = (bufferFromInput = msg => msg.Data, bufferToOutput = buf => buf,
  chunkSize = DEFAULT_CHUNK_SIZE, transformOptions = {}) => {
  let buffer = Buffer.alloc(0)

  const push = (transform, buffer) => {
    if (!transform._readableState.ended) {
      transform.push(bufferToOutput(buffer))
    }
  }

  const transform = new Transform({
    autoDestroy: true,
    transform (msg, encoding, callback) {
      buffer = Buffer.concat([buffer, bufferFromInput(msg)])
      if (buffer.length > chunkSize) {
        push(this, buffer.slice(0, chunkSize))
        buffer = buffer.slice(chunkSize)
      }
      callback()
    },

    flush (callback) {
      push(this, buffer)
      callback()
    },
    ...transformOptions
  })

  transform.on('unpipe', src => {
    // Make sure we cancel any GRPC stream
    if (typeof src.cancel === 'function') {
      src.cancel()
    }
  })
  return transform
}

const BytesToObjects = (bufferToObject = buf => Plaintext(buf), chunkSize = DEFAULT_CHUNK_SIZE) =>
  BufferTransform(msg => msg, bufferToObject, chunkSize,
    { readableObjectMode: true, writableObjectMode: false })

// Converts an object mode stream to a byte mode stream selecting a single buffer from the source using dataSelector
// and buffering bytes in chunks of at chunkSize
const ObjectsToBytes = (bufferSelector = msg => msg.Data, chunkSize = DEFAULT_CHUNK_SIZE) =>
  BufferTransform(bufferSelector, msg => msg, chunkSize,
    { readableObjectMode: false, writableObjectMode: true })

async function write (stream, messages) {
  if (isIterable(messages)) {
    for (const msg of messages) {
      stream.write(msg)
    }
  } else {
    stream.write(messages)
  }
  return stream.close()
}

// Read takes a server streaming function and calls it with the args provided having extracted the reducer and optional
// initial accumulator from the final one or two arguments. The reducer must be supplied. If the final argument is
// a function it is used as the reducer otherwie the final argument is taken to be the initial accumulator value.
// The reducer is a function (acc, val, returnNow) => newAcc where returnNow constructs an early-exiting termination
// value when called like (acc, val, returnNow) => returnNow(newAcc) which causes the read fucntion to return with
// newAcc and cancel() (if available - i.e. GRPC stream) or destroy() the stream.
async function read (stream, ...args) {
  if (args.length < 1) {
    throw new Error('reduce expects at least one reducer function argument')
  }
  let accumulator
  let reducer = args.pop()

  // Try last argument as reducer
  if (typeof reducer !== 'function') {
    if (args.length < 1) {
      throw new Error(`reduce expects at least one callback argument but last argument is ${JSON.stringify(reducer)}`)
    }
    // assume last argument is in fact the initial accumulator value
    accumulator = reducer
    reducer = args.pop()
    // penultimate argument must now be the reducer
    if (typeof reducer !== 'function') {
      throw new Error(`reduce expects a reducer function but penultimate argument is ${JSON.stringify(reducer)}`)
    }
  }

  // Provide function that performs early exit
  const returnNow = ret => {
    // Cancel the stream but defer until next tick to avoid write after end from upstream
    process.nextTick(() => {
      if (typeof stream.cancel === 'function') {
        stream.cancel()
      } else {
        stream.destroy()
      }
    })
    // Scheduling variation in event loops means that we may receive messages after cancel, so prevent them from
    // having any effect by making reducer the identity function
    reducer = acc => acc
    // Return the wrapped value as the final return value of the reducer
    return ret
  }

  return new Promise((resolve, reject) => {
    stream.on('data', data => {
      if (accumulator === undefined) {
        accumulator = data
      } else {
        accumulator = reducer(accumulator, data, returnNow)
      }
    })
    // Allow us (the client) to cancel mid-stream without throwing
    stream.on('error', err => err.code === grpc.status.CANCELLED ? resolve(accumulator) : reject(err))
    stream.on('close', () => resolve(accumulator))
    stream.on('end', () => resolve(accumulator))
  })
}

function reduceLengthPrefixed (byteLength) {
  let length = 0
  return (buffer, data, returnNow) => {
    // Grow the buffer
    buffer = Buffer.concat([buffer, data])
    // First try to read the length prefix
    if (length === 0) {
      if (buffer.length >= byteLength) {
        length = buffer.readUIntBE(0, byteLength)
        buffer = buffer.slice(byteLength)
      }
    }
    if (length > 0 && buffer.length >= length) {
      // If we have read the length prefix and it is contained within
      // our current buffer then we are done
      return returnNow(buffer.slice(0, length))
    }
    // Keep growing the buffer until it contains sufficient data
    return buffer
  }
}

// Reads from the provided byteStream until a varint length-prefixed prefix of the stream can be returned.
// The stream will be destroyed once the prefix has been read so will not be totally consumed
async function readLengthPrefixed (byteStream, byteLength) {
  return read(byteStream, reduceLengthPrefixed(byteLength), Buffer.alloc(0))
}

function base64ify (obj) {
  const newObj = {}
  for (const key of Object.keys(obj)) {
    const value = obj[key]
    if (value instanceof Buffer) {
      newObj[key] = value.toString('base64')
    } else {
      newObj[key] = value
    }
  }
  return newObj
}

module.exports = {
  Client,
  ObjectsToBytes,
  BytesToObjects,
  read,
  write,
  readLengthPrefixed,
  base64ify,
  Plaintext,
  PlaintextAndGrantSpec
}

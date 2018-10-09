// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var hoard_pb = require('./hoard_pb.js');

function serialize_core_Address(arg) {
  if (!(arg instanceof hoard_pb.Address)) {
    throw new Error('Expected argument of type core.Address');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_core_Address(buffer_arg) {
  return hoard_pb.Address.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_core_Ciphertext(arg) {
  if (!(arg instanceof hoard_pb.Ciphertext)) {
    throw new Error('Expected argument of type core.Ciphertext');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_core_Ciphertext(buffer_arg) {
  return hoard_pb.Ciphertext.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_core_Plaintext(arg) {
  if (!(arg instanceof hoard_pb.Plaintext)) {
    throw new Error('Expected argument of type core.Plaintext');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_core_Plaintext(buffer_arg) {
  return hoard_pb.Plaintext.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_core_Reference(arg) {
  if (!(arg instanceof hoard_pb.Reference)) {
    throw new Error('Expected argument of type core.Reference');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_core_Reference(buffer_arg) {
  return hoard_pb.Reference.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_core_ReferenceAndCiphertext(arg) {
  if (!(arg instanceof hoard_pb.ReferenceAndCiphertext)) {
    throw new Error('Expected argument of type core.ReferenceAndCiphertext');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_core_ReferenceAndCiphertext(buffer_arg) {
  return hoard_pb.ReferenceAndCiphertext.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_core_StatInfo(arg) {
  if (!(arg instanceof hoard_pb.StatInfo)) {
    throw new Error('Expected argument of type core.StatInfo');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_core_StatInfo(buffer_arg) {
  return hoard_pb.StatInfo.deserializeBinary(new Uint8Array(buffer_arg));
}


// Provide plaintext and get plaintext back
var CleartextService = exports.CleartextService = {
  // Provide a secret reference to an encrypted blob and get the plaintext
  // data back.
  get: {
    path: '/core.Cleartext/Get',
    requestStream: false,
    responseStream: false,
    requestType: hoard_pb.Reference,
    responseType: hoard_pb.Plaintext,
    requestSerialize: serialize_core_Reference,
    requestDeserialize: deserialize_core_Reference,
    responseSerialize: serialize_core_Plaintext,
    responseDeserialize: deserialize_core_Plaintext,
  },
  // Push some plaintext data into storage and get its deterministically
  // generated secret reference.
  put: {
    path: '/core.Cleartext/Put',
    requestStream: false,
    responseStream: false,
    requestType: hoard_pb.Plaintext,
    responseType: hoard_pb.Reference,
    requestSerialize: serialize_core_Plaintext,
    requestDeserialize: deserialize_core_Plaintext,
    responseSerialize: serialize_core_Reference,
    responseDeserialize: deserialize_core_Reference,
  },
};

exports.CleartextClient = grpc.makeGenericClientConstructor(CleartextService);
// Deterministic encryption
var EncryptionService = exports.EncryptionService = {
  // Encrypt some data and get its deterministically generated
  // secret reference including its address without storing the data.
  encrypt: {
    path: '/core.Encryption/Encrypt',
    requestStream: false,
    responseStream: false,
    requestType: hoard_pb.Plaintext,
    responseType: hoard_pb.ReferenceAndCiphertext,
    requestSerialize: serialize_core_Plaintext,
    requestDeserialize: deserialize_core_Plaintext,
    responseSerialize: serialize_core_ReferenceAndCiphertext,
    responseDeserialize: deserialize_core_ReferenceAndCiphertext,
  },
  // Decrypt the provided data by supplying it alongside its secret
  // reference. The address is not used for decryption and may be omitted.
  decrypt: {
    path: '/core.Encryption/Decrypt',
    requestStream: false,
    responseStream: false,
    requestType: hoard_pb.ReferenceAndCiphertext,
    responseType: hoard_pb.Plaintext,
    requestSerialize: serialize_core_ReferenceAndCiphertext,
    requestDeserialize: deserialize_core_ReferenceAndCiphertext,
    responseSerialize: serialize_core_Plaintext,
    responseDeserialize: deserialize_core_Plaintext,
  },
};

exports.EncryptionClient = grpc.makeGenericClientConstructor(EncryptionService);
// Interact directly with storage backend
var StorageService = exports.StorageService = {
  // Retrieve the (presumably) encrypted data stored at address.
  pull: {
    path: '/core.Storage/Pull',
    requestStream: false,
    responseStream: false,
    requestType: hoard_pb.Address,
    responseType: hoard_pb.Ciphertext,
    requestSerialize: serialize_core_Address,
    requestDeserialize: deserialize_core_Address,
    responseSerialize: serialize_core_Ciphertext,
    responseDeserialize: deserialize_core_Ciphertext,
  },
  // Insert the (presumably) encrypted data provided and get the its address.
  push: {
    path: '/core.Storage/Push',
    requestStream: false,
    responseStream: false,
    requestType: hoard_pb.Ciphertext,
    responseType: hoard_pb.Address,
    requestSerialize: serialize_core_Ciphertext,
    requestDeserialize: deserialize_core_Ciphertext,
    responseSerialize: serialize_core_Address,
    responseDeserialize: deserialize_core_Address,
  },
  // Get some information about the encrypted blob stored at an address,
  // including whether it exists.
  stat: {
    path: '/core.Storage/Stat',
    requestStream: false,
    responseStream: false,
    requestType: hoard_pb.Address,
    responseType: hoard_pb.StatInfo,
    requestSerialize: serialize_core_Address,
    requestDeserialize: deserialize_core_Address,
    responseSerialize: serialize_core_StatInfo,
    responseDeserialize: deserialize_core_StatInfo,
  },
};

exports.StorageClient = grpc.makeGenericClientConstructor(StorageService);

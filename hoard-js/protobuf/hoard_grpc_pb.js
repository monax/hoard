// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var hoard_pb = require('./hoard_pb.js');
var github_com_gogo_protobuf_gogoproto_gogo_pb = require('./github.com/gogo/protobuf/gogoproto/gogo_pb.js');
var reference_pb = require('./reference_pb.js');
var grant_pb = require('./grant_pb.js');
var storage_pb = require('./storage_pb.js');

function serialize_grant_Grant(arg) {
  if (!(arg instanceof grant_pb.Grant)) {
    throw new Error('Expected argument of type grant.Grant');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_grant_Grant(buffer_arg) {
  return grant_pb.Grant.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_hoard_Address(arg) {
  if (!(arg instanceof hoard_pb.Address)) {
    throw new Error('Expected argument of type hoard.Address');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_hoard_Address(buffer_arg) {
  return hoard_pb.Address.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_hoard_Ciphertext(arg) {
  if (!(arg instanceof hoard_pb.Ciphertext)) {
    throw new Error('Expected argument of type hoard.Ciphertext');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_hoard_Ciphertext(buffer_arg) {
  return hoard_pb.Ciphertext.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_hoard_GrantAndGrantSpec(arg) {
  if (!(arg instanceof hoard_pb.GrantAndGrantSpec)) {
    throw new Error('Expected argument of type hoard.GrantAndGrantSpec');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_hoard_GrantAndGrantSpec(buffer_arg) {
  return hoard_pb.GrantAndGrantSpec.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_hoard_Plaintext(arg) {
  if (!(arg instanceof hoard_pb.Plaintext)) {
    throw new Error('Expected argument of type hoard.Plaintext');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_hoard_Plaintext(buffer_arg) {
  return hoard_pb.Plaintext.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_hoard_PlaintextAndGrantSpec(arg) {
  if (!(arg instanceof hoard_pb.PlaintextAndGrantSpec)) {
    throw new Error('Expected argument of type hoard.PlaintextAndGrantSpec');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_hoard_PlaintextAndGrantSpec(buffer_arg) {
  return hoard_pb.PlaintextAndGrantSpec.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_hoard_ReferenceAndCiphertext(arg) {
  if (!(arg instanceof hoard_pb.ReferenceAndCiphertext)) {
    throw new Error('Expected argument of type hoard.ReferenceAndCiphertext');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_hoard_ReferenceAndCiphertext(buffer_arg) {
  return hoard_pb.ReferenceAndCiphertext.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_hoard_ReferenceAndGrantSpec(arg) {
  if (!(arg instanceof hoard_pb.ReferenceAndGrantSpec)) {
    throw new Error('Expected argument of type hoard.ReferenceAndGrantSpec');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_hoard_ReferenceAndGrantSpec(buffer_arg) {
  return hoard_pb.ReferenceAndGrantSpec.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_reference_Ref(arg) {
  if (!(arg instanceof reference_pb.Ref)) {
    throw new Error('Expected argument of type reference.Ref');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_reference_Ref(buffer_arg) {
  return reference_pb.Ref.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_storage_StatInfo(arg) {
  if (!(arg instanceof storage_pb.StatInfo)) {
    throw new Error('Expected argument of type storage.StatInfo');
  }
  return new Buffer(arg.serializeBinary());
}

function deserialize_storage_StatInfo(buffer_arg) {
  return storage_pb.StatInfo.deserializeBinary(new Uint8Array(buffer_arg));
}


var GrantService = exports.GrantService = {
  // Seal a Reference to create a Grant
  seal: {
    path: '/hoard.Grant/Seal',
    requestStream: false,
    responseStream: false,
    requestType: hoard_pb.ReferenceAndGrantSpec,
    responseType: grant_pb.Grant,
    requestSerialize: serialize_hoard_ReferenceAndGrantSpec,
    requestDeserialize: deserialize_hoard_ReferenceAndGrantSpec,
    responseSerialize: serialize_grant_Grant,
    responseDeserialize: deserialize_grant_Grant,
  },
  // Unseal a Grant to recover the Reference
  unseal: {
    path: '/hoard.Grant/Unseal',
    requestStream: false,
    responseStream: false,
    requestType: grant_pb.Grant,
    responseType: reference_pb.Ref,
    requestSerialize: serialize_grant_Grant,
    requestDeserialize: deserialize_grant_Grant,
    responseSerialize: serialize_reference_Ref,
    responseDeserialize: deserialize_reference_Ref,
  },
  // Convert one grant to another grant to re-share with another party or just
  // to change grant type
  reseal: {
    path: '/hoard.Grant/Reseal',
    requestStream: false,
    responseStream: false,
    requestType: hoard_pb.GrantAndGrantSpec,
    responseType: grant_pb.Grant,
    requestSerialize: serialize_hoard_GrantAndGrantSpec,
    requestDeserialize: deserialize_hoard_GrantAndGrantSpec,
    responseSerialize: serialize_grant_Grant,
    responseDeserialize: deserialize_grant_Grant,
  },
  // Put a Plaintext and returned the sealed Reference as a Grant
  putSeal: {
    path: '/hoard.Grant/PutSeal',
    requestStream: false,
    responseStream: false,
    requestType: hoard_pb.PlaintextAndGrantSpec,
    responseType: grant_pb.Grant,
    requestSerialize: serialize_hoard_PlaintextAndGrantSpec,
    requestDeserialize: deserialize_hoard_PlaintextAndGrantSpec,
    responseSerialize: serialize_grant_Grant,
    responseDeserialize: deserialize_grant_Grant,
  },
  // Unseal a Grant and follow the Reference to return a Plaintext
  unsealGet: {
    path: '/hoard.Grant/UnsealGet',
    requestStream: false,
    responseStream: false,
    requestType: grant_pb.Grant,
    responseType: hoard_pb.Plaintext,
    requestSerialize: serialize_grant_Grant,
    requestDeserialize: deserialize_grant_Grant,
    responseSerialize: serialize_hoard_Plaintext,
    responseDeserialize: deserialize_hoard_Plaintext,
  },
};

exports.GrantClient = grpc.makeGenericClientConstructor(GrantService);
// Provide plaintext and get plaintext back
var CleartextService = exports.CleartextService = {
  // Push some plaintext data into storage and get its deterministically
  // generated secret reference.
  put: {
    path: '/hoard.Cleartext/Put',
    requestStream: false,
    responseStream: false,
    requestType: hoard_pb.Plaintext,
    responseType: reference_pb.Ref,
    requestSerialize: serialize_hoard_Plaintext,
    requestDeserialize: deserialize_hoard_Plaintext,
    responseSerialize: serialize_reference_Ref,
    responseDeserialize: deserialize_reference_Ref,
  },
  // Provide a secret reference to an encrypted blob and get the plaintext
  // data back.
  get: {
    path: '/hoard.Cleartext/Get',
    requestStream: false,
    responseStream: false,
    requestType: reference_pb.Ref,
    responseType: hoard_pb.Plaintext,
    requestSerialize: serialize_reference_Ref,
    requestDeserialize: deserialize_reference_Ref,
    responseSerialize: serialize_hoard_Plaintext,
    responseDeserialize: deserialize_hoard_Plaintext,
  },
};

exports.CleartextClient = grpc.makeGenericClientConstructor(CleartextService);
// Deterministic encryption
var EncryptionService = exports.EncryptionService = {
  // Encrypt some data and get its deterministically generated
  // secret reference including its address without storing the data.
  encrypt: {
    path: '/hoard.Encryption/Encrypt',
    requestStream: false,
    responseStream: false,
    requestType: hoard_pb.Plaintext,
    responseType: hoard_pb.ReferenceAndCiphertext,
    requestSerialize: serialize_hoard_Plaintext,
    requestDeserialize: deserialize_hoard_Plaintext,
    responseSerialize: serialize_hoard_ReferenceAndCiphertext,
    responseDeserialize: deserialize_hoard_ReferenceAndCiphertext,
  },
  // Decrypt the provided data by supplying it alongside its secret
  // reference. The address is not used for decryption and may be omitted.
  decrypt: {
    path: '/hoard.Encryption/Decrypt',
    requestStream: false,
    responseStream: false,
    requestType: hoard_pb.ReferenceAndCiphertext,
    responseType: hoard_pb.Plaintext,
    requestSerialize: serialize_hoard_ReferenceAndCiphertext,
    requestDeserialize: deserialize_hoard_ReferenceAndCiphertext,
    responseSerialize: serialize_hoard_Plaintext,
    responseDeserialize: deserialize_hoard_Plaintext,
  },
};

exports.EncryptionClient = grpc.makeGenericClientConstructor(EncryptionService);
// Interact directly with storage backend
var StorageService = exports.StorageService = {
  // Insert the (presumably) encrypted data provided and get the its address.
  push: {
    path: '/hoard.Storage/Push',
    requestStream: false,
    responseStream: false,
    requestType: hoard_pb.Ciphertext,
    responseType: hoard_pb.Address,
    requestSerialize: serialize_hoard_Ciphertext,
    requestDeserialize: deserialize_hoard_Ciphertext,
    responseSerialize: serialize_hoard_Address,
    responseDeserialize: deserialize_hoard_Address,
  },
  // Retrieve the (presumably) encrypted data stored at address.
  pull: {
    path: '/hoard.Storage/Pull',
    requestStream: false,
    responseStream: false,
    requestType: hoard_pb.Address,
    responseType: hoard_pb.Ciphertext,
    requestSerialize: serialize_hoard_Address,
    requestDeserialize: deserialize_hoard_Address,
    responseSerialize: serialize_hoard_Ciphertext,
    responseDeserialize: deserialize_hoard_Ciphertext,
  },
  // Get some information about the encrypted blob stored at an address,
  // including whether it exists.
  stat: {
    path: '/hoard.Storage/Stat',
    requestStream: false,
    responseStream: false,
    requestType: hoard_pb.Address,
    responseType: storage_pb.StatInfo,
    requestSerialize: serialize_hoard_Address,
    requestDeserialize: deserialize_hoard_Address,
    responseSerialize: serialize_storage_StatInfo,
    responseDeserialize: deserialize_storage_StatInfo,
  },
};

exports.StorageClient = grpc.makeGenericClientConstructor(StorageService);

// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('@grpc/grpc-js');
var api_pb = require('./api_pb.js');
var grant_pb = require('./grant_pb.js');
var reference_pb = require('./reference_pb.js');
var stores_pb = require('./stores_pb.js');

function serialize_api_Address(arg) {
  if (!(arg instanceof api_pb.Address)) {
    throw new Error('Expected argument of type api.Address');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_api_Address(buffer_arg) {
  return api_pb.Address.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_api_Ciphertext(arg) {
  if (!(arg instanceof api_pb.Ciphertext)) {
    throw new Error('Expected argument of type api.Ciphertext');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_api_Ciphertext(buffer_arg) {
  return api_pb.Ciphertext.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_api_GrantAndGrantSpec(arg) {
  if (!(arg instanceof api_pb.GrantAndGrantSpec)) {
    throw new Error('Expected argument of type api.GrantAndGrantSpec');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_api_GrantAndGrantSpec(buffer_arg) {
  return api_pb.GrantAndGrantSpec.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_api_Plaintext(arg) {
  if (!(arg instanceof api_pb.Plaintext)) {
    throw new Error('Expected argument of type api.Plaintext');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_api_Plaintext(buffer_arg) {
  return api_pb.Plaintext.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_api_PlaintextAndGrantSpec(arg) {
  if (!(arg instanceof api_pb.PlaintextAndGrantSpec)) {
    throw new Error('Expected argument of type api.PlaintextAndGrantSpec');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_api_PlaintextAndGrantSpec(buffer_arg) {
  return api_pb.PlaintextAndGrantSpec.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_api_ReferenceAndCiphertext(arg) {
  if (!(arg instanceof api_pb.ReferenceAndCiphertext)) {
    throw new Error('Expected argument of type api.ReferenceAndCiphertext');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_api_ReferenceAndCiphertext(buffer_arg) {
  return api_pb.ReferenceAndCiphertext.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_api_ReferenceAndGrantSpec(arg) {
  if (!(arg instanceof api_pb.ReferenceAndGrantSpec)) {
    throw new Error('Expected argument of type api.ReferenceAndGrantSpec');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_api_ReferenceAndGrantSpec(buffer_arg) {
  return api_pb.ReferenceAndGrantSpec.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_grant_Grant(arg) {
  if (!(arg instanceof grant_pb.Grant)) {
    throw new Error('Expected argument of type grant.Grant');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_grant_Grant(buffer_arg) {
  return grant_pb.Grant.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_reference_Ref(arg) {
  if (!(arg instanceof reference_pb.Ref)) {
    throw new Error('Expected argument of type reference.Ref');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_reference_Ref(buffer_arg) {
  return reference_pb.Ref.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_stores_StatInfo(arg) {
  if (!(arg instanceof stores_pb.StatInfo)) {
    throw new Error('Expected argument of type stores.StatInfo');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_stores_StatInfo(buffer_arg) {
  return stores_pb.StatInfo.deserializeBinary(new Uint8Array(buffer_arg));
}


var GrantService = exports.GrantService = {
  // Put a Plaintext and returned the sealed Reference as a Grant
putSeal: {
    path: '/api.Grant/PutSeal',
    requestStream: true,
    responseStream: false,
    requestType: api_pb.PlaintextAndGrantSpec,
    responseType: grant_pb.Grant,
    requestSerialize: serialize_api_PlaintextAndGrantSpec,
    requestDeserialize: deserialize_api_PlaintextAndGrantSpec,
    responseSerialize: serialize_grant_Grant,
    responseDeserialize: deserialize_grant_Grant,
  },
  // Unseal a Grant and follow the Reference to return a Plaintext
unsealGet: {
    path: '/api.Grant/UnsealGet',
    requestStream: false,
    responseStream: true,
    requestType: grant_pb.Grant,
    responseType: api_pb.Plaintext,
    requestSerialize: serialize_grant_Grant,
    requestDeserialize: deserialize_grant_Grant,
    responseSerialize: serialize_api_Plaintext,
    responseDeserialize: deserialize_api_Plaintext,
  },
  // Seal a Reference to create a Grant
seal: {
    path: '/api.Grant/Seal',
    requestStream: true,
    responseStream: false,
    requestType: api_pb.ReferenceAndGrantSpec,
    responseType: grant_pb.Grant,
    requestSerialize: serialize_api_ReferenceAndGrantSpec,
    requestDeserialize: deserialize_api_ReferenceAndGrantSpec,
    responseSerialize: serialize_grant_Grant,
    responseDeserialize: deserialize_grant_Grant,
  },
  // Unseal a Grant to recover the Reference
unseal: {
    path: '/api.Grant/Unseal',
    requestStream: false,
    responseStream: true,
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
    path: '/api.Grant/Reseal',
    requestStream: false,
    responseStream: false,
    requestType: api_pb.GrantAndGrantSpec,
    responseType: grant_pb.Grant,
    requestSerialize: serialize_api_GrantAndGrantSpec,
    requestDeserialize: deserialize_api_GrantAndGrantSpec,
    responseSerialize: serialize_grant_Grant,
    responseDeserialize: deserialize_grant_Grant,
  },
  // Unseal a Grant and follow the Reference to delete the Plaintext
unsealDelete: {
    path: '/api.Grant/UnsealDelete',
    requestStream: false,
    responseStream: true,
    requestType: grant_pb.Grant,
    responseType: api_pb.Address,
    requestSerialize: serialize_grant_Grant,
    requestDeserialize: deserialize_grant_Grant,
    responseSerialize: serialize_api_Address,
    responseDeserialize: deserialize_api_Address,
  },
};

exports.GrantClient = grpc.makeGenericClientConstructor(GrantService);
// Provide plaintext and get plaintext back
var CleartextService = exports.CleartextService = {
  // Push some plaintext data into storage and get its deterministically
// generated secret reference.
put: {
    path: '/api.Cleartext/Put',
    requestStream: true,
    responseStream: true,
    requestType: api_pb.Plaintext,
    responseType: reference_pb.Ref,
    requestSerialize: serialize_api_Plaintext,
    requestDeserialize: deserialize_api_Plaintext,
    responseSerialize: serialize_reference_Ref,
    responseDeserialize: deserialize_reference_Ref,
  },
  // Provide a secret reference to an encrypted blob and get the plaintext
// data back.
get: {
    path: '/api.Cleartext/Get',
    requestStream: true,
    responseStream: true,
    requestType: reference_pb.Ref,
    responseType: api_pb.Plaintext,
    requestSerialize: serialize_reference_Ref,
    requestDeserialize: deserialize_reference_Ref,
    responseSerialize: serialize_api_Plaintext,
    responseDeserialize: deserialize_api_Plaintext,
  },
};

exports.CleartextClient = grpc.makeGenericClientConstructor(CleartextService);
// Deterministic encryption
var EncryptionService = exports.EncryptionService = {
  // Encrypt some data and get its deterministically generated
// secret reference including its address without storing the data.
encrypt: {
    path: '/api.Encryption/Encrypt',
    requestStream: true,
    responseStream: true,
    requestType: api_pb.Plaintext,
    responseType: api_pb.ReferenceAndCiphertext,
    requestSerialize: serialize_api_Plaintext,
    requestDeserialize: deserialize_api_Plaintext,
    responseSerialize: serialize_api_ReferenceAndCiphertext,
    responseDeserialize: deserialize_api_ReferenceAndCiphertext,
  },
  // Decrypt the provided data by supplying it alongside its secret
// reference. The address is not used for decryption and may be omitted.
decrypt: {
    path: '/api.Encryption/Decrypt',
    requestStream: true,
    responseStream: true,
    requestType: api_pb.ReferenceAndCiphertext,
    responseType: api_pb.Plaintext,
    requestSerialize: serialize_api_ReferenceAndCiphertext,
    requestDeserialize: deserialize_api_ReferenceAndCiphertext,
    responseSerialize: serialize_api_Plaintext,
    responseDeserialize: deserialize_api_Plaintext,
  },
};

exports.EncryptionClient = grpc.makeGenericClientConstructor(EncryptionService);
// Interact directly with storage backend
var StorageService = exports.StorageService = {
  // Insert the (presumably) encrypted data provided and get the its address.
push: {
    path: '/api.Storage/Push',
    requestStream: true,
    responseStream: true,
    requestType: api_pb.Ciphertext,
    responseType: api_pb.Address,
    requestSerialize: serialize_api_Ciphertext,
    requestDeserialize: deserialize_api_Ciphertext,
    responseSerialize: serialize_api_Address,
    responseDeserialize: deserialize_api_Address,
  },
  // Retrieve the (presumably) encrypted data stored at address.
pull: {
    path: '/api.Storage/Pull',
    requestStream: true,
    responseStream: true,
    requestType: api_pb.Address,
    responseType: api_pb.Ciphertext,
    requestSerialize: serialize_api_Address,
    requestDeserialize: deserialize_api_Address,
    responseSerialize: serialize_api_Ciphertext,
    responseDeserialize: deserialize_api_Ciphertext,
  },
  // Get some information about the encrypted blob stored at an address,
// including whether it exists.
stat: {
    path: '/api.Storage/Stat',
    requestStream: false,
    responseStream: false,
    requestType: api_pb.Address,
    responseType: stores_pb.StatInfo,
    requestSerialize: serialize_api_Address,
    requestDeserialize: deserialize_api_Address,
    responseSerialize: serialize_stores_StatInfo,
    responseDeserialize: deserialize_stores_StatInfo,
  },
  // Delete the encrypted blob stored at address
delete: {
    path: '/api.Storage/Delete',
    requestStream: false,
    responseStream: false,
    requestType: api_pb.Address,
    responseType: api_pb.Address,
    requestSerialize: serialize_api_Address,
    requestDeserialize: deserialize_api_Address,
    responseSerialize: serialize_api_Address,
    responseDeserialize: deserialize_api_Address,
  },
};

exports.StorageClient = grpc.makeGenericClientConstructor(StorageService);

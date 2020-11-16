// GENERATED CODE -- DO NOT EDIT!

// package: api
// file: api.proto

import * as api_pb from "./api_pb";
import * as grant_pb from "./grant_pb";
import * as reference_pb from "./reference_pb";
import * as stores_pb from "./stores_pb";
import * as grpc from "@grpc/grpc-js";

interface IGrantService extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
  seal: grpc.MethodDefinition<api_pb.ReferenceAndGrantSpec, grant_pb.Grant>;
  unseal: grpc.MethodDefinition<grant_pb.Grant, reference_pb.Ref>;
  reseal: grpc.MethodDefinition<api_pb.GrantAndGrantSpec, grant_pb.Grant>;
  putSeal: grpc.MethodDefinition<api_pb.PlaintextAndGrantSpec, grant_pb.Grant>;
  unsealGet: grpc.MethodDefinition<grant_pb.Grant, api_pb.Plaintext>;
  unsealDelete: grpc.MethodDefinition<grant_pb.Grant, api_pb.Address>;
}

export const GrantService: IGrantService;

export class GrantClient extends grpc.Client {
  constructor(address: string, credentials: grpc.ChannelCredentials, options?: object);
  seal(callback: grpc.requestCallback<grant_pb.Grant>): grpc.ClientWritableStream<api_pb.ReferenceAndGrantSpec>;
  seal(metadataOrOptions: grpc.Metadata | grpc.CallOptions | null, callback: grpc.requestCallback<grant_pb.Grant>): grpc.ClientWritableStream<api_pb.ReferenceAndGrantSpec>;
  seal(metadata: grpc.Metadata | null, options: grpc.CallOptions | null, callback: grpc.requestCallback<grant_pb.Grant>): grpc.ClientWritableStream<api_pb.ReferenceAndGrantSpec>;
  unseal(argument: grant_pb.Grant, metadataOrOptions?: grpc.Metadata | grpc.CallOptions | null): grpc.ClientReadableStream<reference_pb.Ref>;
  unseal(argument: grant_pb.Grant, metadata?: grpc.Metadata | null, options?: grpc.CallOptions | null): grpc.ClientReadableStream<reference_pb.Ref>;
  reseal(argument: api_pb.GrantAndGrantSpec, callback: grpc.requestCallback<grant_pb.Grant>): grpc.ClientUnaryCall;
  reseal(argument: api_pb.GrantAndGrantSpec, metadataOrOptions: grpc.Metadata | grpc.CallOptions | null, callback: grpc.requestCallback<grant_pb.Grant>): grpc.ClientUnaryCall;
  reseal(argument: api_pb.GrantAndGrantSpec, metadata: grpc.Metadata | null, options: grpc.CallOptions | null, callback: grpc.requestCallback<grant_pb.Grant>): grpc.ClientUnaryCall;
  putSeal(callback: grpc.requestCallback<grant_pb.Grant>): grpc.ClientWritableStream<api_pb.PlaintextAndGrantSpec>;
  putSeal(metadataOrOptions: grpc.Metadata | grpc.CallOptions | null, callback: grpc.requestCallback<grant_pb.Grant>): grpc.ClientWritableStream<api_pb.PlaintextAndGrantSpec>;
  putSeal(metadata: grpc.Metadata | null, options: grpc.CallOptions | null, callback: grpc.requestCallback<grant_pb.Grant>): grpc.ClientWritableStream<api_pb.PlaintextAndGrantSpec>;
  unsealGet(argument: grant_pb.Grant, metadataOrOptions?: grpc.Metadata | grpc.CallOptions | null): grpc.ClientReadableStream<api_pb.Plaintext>;
  unsealGet(argument: grant_pb.Grant, metadata?: grpc.Metadata | null, options?: grpc.CallOptions | null): grpc.ClientReadableStream<api_pb.Plaintext>;
  unsealDelete(argument: grant_pb.Grant, metadataOrOptions?: grpc.Metadata | grpc.CallOptions | null): grpc.ClientReadableStream<api_pb.Address>;
  unsealDelete(argument: grant_pb.Grant, metadata?: grpc.Metadata | null, options?: grpc.CallOptions | null): grpc.ClientReadableStream<api_pb.Address>;
}

interface ICleartextService extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
  put: grpc.MethodDefinition<api_pb.Plaintext, reference_pb.Ref>;
  get: grpc.MethodDefinition<reference_pb.Ref, api_pb.Plaintext>;
}

export const CleartextService: ICleartextService;

export class CleartextClient extends grpc.Client {
  constructor(address: string, credentials: grpc.ChannelCredentials, options?: object);
  put(metadataOrOptions?: grpc.Metadata | grpc.CallOptions | null): grpc.ClientDuplexStream<api_pb.Plaintext, reference_pb.Ref>;
  put(metadata?: grpc.Metadata | null, options?: grpc.CallOptions | null): grpc.ClientDuplexStream<api_pb.Plaintext, reference_pb.Ref>;
  get(metadataOrOptions?: grpc.Metadata | grpc.CallOptions | null): grpc.ClientDuplexStream<reference_pb.Ref, api_pb.Plaintext>;
  get(metadata?: grpc.Metadata | null, options?: grpc.CallOptions | null): grpc.ClientDuplexStream<reference_pb.Ref, api_pb.Plaintext>;
}

interface IEncryptionService extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
  encrypt: grpc.MethodDefinition<api_pb.Plaintext, api_pb.ReferenceAndCiphertext>;
  decrypt: grpc.MethodDefinition<api_pb.ReferenceAndCiphertext, api_pb.Plaintext>;
}

export const EncryptionService: IEncryptionService;

export class EncryptionClient extends grpc.Client {
  constructor(address: string, credentials: grpc.ChannelCredentials, options?: object);
  encrypt(metadataOrOptions?: grpc.Metadata | grpc.CallOptions | null): grpc.ClientDuplexStream<api_pb.Plaintext, api_pb.ReferenceAndCiphertext>;
  encrypt(metadata?: grpc.Metadata | null, options?: grpc.CallOptions | null): grpc.ClientDuplexStream<api_pb.Plaintext, api_pb.ReferenceAndCiphertext>;
  decrypt(metadataOrOptions?: grpc.Metadata | grpc.CallOptions | null): grpc.ClientDuplexStream<api_pb.ReferenceAndCiphertext, api_pb.Plaintext>;
  decrypt(metadata?: grpc.Metadata | null, options?: grpc.CallOptions | null): grpc.ClientDuplexStream<api_pb.ReferenceAndCiphertext, api_pb.Plaintext>;
}

interface IStorageService extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
  push: grpc.MethodDefinition<api_pb.Ciphertext, api_pb.Address>;
  pull: grpc.MethodDefinition<api_pb.Address, api_pb.Ciphertext>;
  stat: grpc.MethodDefinition<api_pb.Address, stores_pb.StatInfo>;
  delete: grpc.MethodDefinition<api_pb.Address, api_pb.Address>;
}

export const StorageService: IStorageService;

export class StorageClient extends grpc.Client {
  constructor(address: string, credentials: grpc.ChannelCredentials, options?: object);
  push(metadataOrOptions?: grpc.Metadata | grpc.CallOptions | null): grpc.ClientDuplexStream<api_pb.Ciphertext, api_pb.Address>;
  push(metadata?: grpc.Metadata | null, options?: grpc.CallOptions | null): grpc.ClientDuplexStream<api_pb.Ciphertext, api_pb.Address>;
  pull(metadataOrOptions?: grpc.Metadata | grpc.CallOptions | null): grpc.ClientDuplexStream<api_pb.Address, api_pb.Ciphertext>;
  pull(metadata?: grpc.Metadata | null, options?: grpc.CallOptions | null): grpc.ClientDuplexStream<api_pb.Address, api_pb.Ciphertext>;
  stat(argument: api_pb.Address, callback: grpc.requestCallback<stores_pb.StatInfo>): grpc.ClientUnaryCall;
  stat(argument: api_pb.Address, metadataOrOptions: grpc.Metadata | grpc.CallOptions | null, callback: grpc.requestCallback<stores_pb.StatInfo>): grpc.ClientUnaryCall;
  stat(argument: api_pb.Address, metadata: grpc.Metadata | null, options: grpc.CallOptions | null, callback: grpc.requestCallback<stores_pb.StatInfo>): grpc.ClientUnaryCall;
  delete(argument: api_pb.Address, callback: grpc.requestCallback<api_pb.Address>): grpc.ClientUnaryCall;
  delete(argument: api_pb.Address, metadataOrOptions: grpc.Metadata | grpc.CallOptions | null, callback: grpc.requestCallback<api_pb.Address>): grpc.ClientUnaryCall;
  delete(argument: api_pb.Address, metadata: grpc.Metadata | null, options: grpc.CallOptions | null, callback: grpc.requestCallback<api_pb.Address>): grpc.ClientUnaryCall;
}

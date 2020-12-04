// package: api
// file: api.proto

/* tslint:disable */
/* eslint-disable */

import * as grpc from "@grpc/grpc-js";
import {handleClientStreamingCall} from "@grpc/grpc-js/build/src/server-call";
import * as api_pb from "./api_pb";
import * as grant_pb from "./grant_pb";
import * as reference_pb from "./reference_pb";
import * as stores_pb from "./stores_pb";

interface IGrantService extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
    putSeal: IGrantService_IPutSeal;
    unsealGet: IGrantService_IUnsealGet;
    seal: IGrantService_ISeal;
    unseal: IGrantService_IUnseal;
    reseal: IGrantService_IReseal;
    unsealDelete: IGrantService_IUnsealDelete;
}

interface IGrantService_IPutSeal extends grpc.MethodDefinition<api_pb.PlaintextAndGrantSpec, grant_pb.Grant> {
    path: "/api.Grant/PutSeal";
    requestStream: true;
    responseStream: false;
    requestSerialize: grpc.serialize<api_pb.PlaintextAndGrantSpec>;
    requestDeserialize: grpc.deserialize<api_pb.PlaintextAndGrantSpec>;
    responseSerialize: grpc.serialize<grant_pb.Grant>;
    responseDeserialize: grpc.deserialize<grant_pb.Grant>;
}
interface IGrantService_IUnsealGet extends grpc.MethodDefinition<grant_pb.Grant, api_pb.Plaintext> {
    path: "/api.Grant/UnsealGet";
    requestStream: false;
    responseStream: true;
    requestSerialize: grpc.serialize<grant_pb.Grant>;
    requestDeserialize: grpc.deserialize<grant_pb.Grant>;
    responseSerialize: grpc.serialize<api_pb.Plaintext>;
    responseDeserialize: grpc.deserialize<api_pb.Plaintext>;
}
interface IGrantService_ISeal extends grpc.MethodDefinition<api_pb.ReferenceAndGrantSpec, grant_pb.Grant> {
    path: "/api.Grant/Seal";
    requestStream: true;
    responseStream: false;
    requestSerialize: grpc.serialize<api_pb.ReferenceAndGrantSpec>;
    requestDeserialize: grpc.deserialize<api_pb.ReferenceAndGrantSpec>;
    responseSerialize: grpc.serialize<grant_pb.Grant>;
    responseDeserialize: grpc.deserialize<grant_pb.Grant>;
}
interface IGrantService_IUnseal extends grpc.MethodDefinition<grant_pb.Grant, reference_pb.Ref> {
    path: "/api.Grant/Unseal";
    requestStream: false;
    responseStream: true;
    requestSerialize: grpc.serialize<grant_pb.Grant>;
    requestDeserialize: grpc.deserialize<grant_pb.Grant>;
    responseSerialize: grpc.serialize<reference_pb.Ref>;
    responseDeserialize: grpc.deserialize<reference_pb.Ref>;
}
interface IGrantService_IReseal extends grpc.MethodDefinition<api_pb.GrantAndGrantSpec, grant_pb.Grant> {
    path: "/api.Grant/Reseal";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<api_pb.GrantAndGrantSpec>;
    requestDeserialize: grpc.deserialize<api_pb.GrantAndGrantSpec>;
    responseSerialize: grpc.serialize<grant_pb.Grant>;
    responseDeserialize: grpc.deserialize<grant_pb.Grant>;
}
interface IGrantService_IUnsealDelete extends grpc.MethodDefinition<grant_pb.Grant, api_pb.Address> {
    path: "/api.Grant/UnsealDelete";
    requestStream: false;
    responseStream: true;
    requestSerialize: grpc.serialize<grant_pb.Grant>;
    requestDeserialize: grpc.deserialize<grant_pb.Grant>;
    responseSerialize: grpc.serialize<api_pb.Address>;
    responseDeserialize: grpc.deserialize<api_pb.Address>;
}

export const GrantService: IGrantService;

export interface IGrantServer {
    putSeal: handleClientStreamingCall<api_pb.PlaintextAndGrantSpec, grant_pb.Grant>;
    unsealGet: grpc.handleServerStreamingCall<grant_pb.Grant, api_pb.Plaintext>;
    seal: handleClientStreamingCall<api_pb.ReferenceAndGrantSpec, grant_pb.Grant>;
    unseal: grpc.handleServerStreamingCall<grant_pb.Grant, reference_pb.Ref>;
    reseal: grpc.handleUnaryCall<api_pb.GrantAndGrantSpec, grant_pb.Grant>;
    unsealDelete: grpc.handleServerStreamingCall<grant_pb.Grant, api_pb.Address>;
}

export interface IGrantClient {
    putSeal(callback: (error: grpc.ServiceError | null, response: grant_pb.Grant) => void): grpc.ClientWritableStream<api_pb.PlaintextAndGrantSpec>;
    putSeal(metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: grant_pb.Grant) => void): grpc.ClientWritableStream<api_pb.PlaintextAndGrantSpec>;
    putSeal(options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: grant_pb.Grant) => void): grpc.ClientWritableStream<api_pb.PlaintextAndGrantSpec>;
    putSeal(metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: grant_pb.Grant) => void): grpc.ClientWritableStream<api_pb.PlaintextAndGrantSpec>;
    unsealGet(request: grant_pb.Grant, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<api_pb.Plaintext>;
    unsealGet(request: grant_pb.Grant, metadata?: grpc.Metadata, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<api_pb.Plaintext>;
    seal(callback: (error: grpc.ServiceError | null, response: grant_pb.Grant) => void): grpc.ClientWritableStream<api_pb.ReferenceAndGrantSpec>;
    seal(metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: grant_pb.Grant) => void): grpc.ClientWritableStream<api_pb.ReferenceAndGrantSpec>;
    seal(options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: grant_pb.Grant) => void): grpc.ClientWritableStream<api_pb.ReferenceAndGrantSpec>;
    seal(metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: grant_pb.Grant) => void): grpc.ClientWritableStream<api_pb.ReferenceAndGrantSpec>;
    unseal(request: grant_pb.Grant, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<reference_pb.Ref>;
    unseal(request: grant_pb.Grant, metadata?: grpc.Metadata, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<reference_pb.Ref>;
    reseal(request: api_pb.GrantAndGrantSpec, callback: (error: grpc.ServiceError | null, response: grant_pb.Grant) => void): grpc.ClientUnaryCall;
    reseal(request: api_pb.GrantAndGrantSpec, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: grant_pb.Grant) => void): grpc.ClientUnaryCall;
    reseal(request: api_pb.GrantAndGrantSpec, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: grant_pb.Grant) => void): grpc.ClientUnaryCall;
    unsealDelete(request: grant_pb.Grant, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<api_pb.Address>;
    unsealDelete(request: grant_pb.Grant, metadata?: grpc.Metadata, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<api_pb.Address>;
}

export class GrantClient extends grpc.Client implements IGrantClient {
    constructor(address: string, credentials: grpc.ChannelCredentials, options?: Partial<grpc.ClientOptions>);
    public putSeal(callback: (error: grpc.ServiceError | null, response: grant_pb.Grant) => void): grpc.ClientWritableStream<api_pb.PlaintextAndGrantSpec>;
    public putSeal(metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: grant_pb.Grant) => void): grpc.ClientWritableStream<api_pb.PlaintextAndGrantSpec>;
    public putSeal(options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: grant_pb.Grant) => void): grpc.ClientWritableStream<api_pb.PlaintextAndGrantSpec>;
    public putSeal(metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: grant_pb.Grant) => void): grpc.ClientWritableStream<api_pb.PlaintextAndGrantSpec>;
    public unsealGet(request: grant_pb.Grant, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<api_pb.Plaintext>;
    public unsealGet(request: grant_pb.Grant, metadata?: grpc.Metadata, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<api_pb.Plaintext>;
    public seal(callback: (error: grpc.ServiceError | null, response: grant_pb.Grant) => void): grpc.ClientWritableStream<api_pb.ReferenceAndGrantSpec>;
    public seal(metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: grant_pb.Grant) => void): grpc.ClientWritableStream<api_pb.ReferenceAndGrantSpec>;
    public seal(options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: grant_pb.Grant) => void): grpc.ClientWritableStream<api_pb.ReferenceAndGrantSpec>;
    public seal(metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: grant_pb.Grant) => void): grpc.ClientWritableStream<api_pb.ReferenceAndGrantSpec>;
    public unseal(request: grant_pb.Grant, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<reference_pb.Ref>;
    public unseal(request: grant_pb.Grant, metadata?: grpc.Metadata, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<reference_pb.Ref>;
    public reseal(request: api_pb.GrantAndGrantSpec, callback: (error: grpc.ServiceError | null, response: grant_pb.Grant) => void): grpc.ClientUnaryCall;
    public reseal(request: api_pb.GrantAndGrantSpec, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: grant_pb.Grant) => void): grpc.ClientUnaryCall;
    public reseal(request: api_pb.GrantAndGrantSpec, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: grant_pb.Grant) => void): grpc.ClientUnaryCall;
    public unsealDelete(request: grant_pb.Grant, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<api_pb.Address>;
    public unsealDelete(request: grant_pb.Grant, metadata?: grpc.Metadata, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<api_pb.Address>;
}

interface ICleartextService extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
    put: ICleartextService_IPut;
    get: ICleartextService_IGet;
}

interface ICleartextService_IPut extends grpc.MethodDefinition<api_pb.Plaintext, reference_pb.Ref> {
    path: "/api.Cleartext/Put";
    requestStream: true;
    responseStream: true;
    requestSerialize: grpc.serialize<api_pb.Plaintext>;
    requestDeserialize: grpc.deserialize<api_pb.Plaintext>;
    responseSerialize: grpc.serialize<reference_pb.Ref>;
    responseDeserialize: grpc.deserialize<reference_pb.Ref>;
}
interface ICleartextService_IGet extends grpc.MethodDefinition<reference_pb.Ref, api_pb.Plaintext> {
    path: "/api.Cleartext/Get";
    requestStream: true;
    responseStream: true;
    requestSerialize: grpc.serialize<reference_pb.Ref>;
    requestDeserialize: grpc.deserialize<reference_pb.Ref>;
    responseSerialize: grpc.serialize<api_pb.Plaintext>;
    responseDeserialize: grpc.deserialize<api_pb.Plaintext>;
}

export const CleartextService: ICleartextService;

export interface ICleartextServer {
    put: grpc.handleBidiStreamingCall<api_pb.Plaintext, reference_pb.Ref>;
    get: grpc.handleBidiStreamingCall<reference_pb.Ref, api_pb.Plaintext>;
}

export interface ICleartextClient {
    put(): grpc.ClientDuplexStream<api_pb.Plaintext, reference_pb.Ref>;
    put(options: Partial<grpc.CallOptions>): grpc.ClientDuplexStream<api_pb.Plaintext, reference_pb.Ref>;
    put(metadata: grpc.Metadata, options?: Partial<grpc.CallOptions>): grpc.ClientDuplexStream<api_pb.Plaintext, reference_pb.Ref>;
    get(): grpc.ClientDuplexStream<reference_pb.Ref, api_pb.Plaintext>;
    get(options: Partial<grpc.CallOptions>): grpc.ClientDuplexStream<reference_pb.Ref, api_pb.Plaintext>;
    get(metadata: grpc.Metadata, options?: Partial<grpc.CallOptions>): grpc.ClientDuplexStream<reference_pb.Ref, api_pb.Plaintext>;
}

export class CleartextClient extends grpc.Client implements ICleartextClient {
    constructor(address: string, credentials: grpc.ChannelCredentials, options?: Partial<grpc.ClientOptions>);
    public put(options?: Partial<grpc.CallOptions>): grpc.ClientDuplexStream<api_pb.Plaintext, reference_pb.Ref>;
    public put(metadata?: grpc.Metadata, options?: Partial<grpc.CallOptions>): grpc.ClientDuplexStream<api_pb.Plaintext, reference_pb.Ref>;
    public get(options?: Partial<grpc.CallOptions>): grpc.ClientDuplexStream<reference_pb.Ref, api_pb.Plaintext>;
    public get(metadata?: grpc.Metadata, options?: Partial<grpc.CallOptions>): grpc.ClientDuplexStream<reference_pb.Ref, api_pb.Plaintext>;
}

interface IEncryptionService extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
    encrypt: IEncryptionService_IEncrypt;
    decrypt: IEncryptionService_IDecrypt;
}

interface IEncryptionService_IEncrypt extends grpc.MethodDefinition<api_pb.Plaintext, api_pb.ReferenceAndCiphertext> {
    path: "/api.Encryption/Encrypt";
    requestStream: true;
    responseStream: true;
    requestSerialize: grpc.serialize<api_pb.Plaintext>;
    requestDeserialize: grpc.deserialize<api_pb.Plaintext>;
    responseSerialize: grpc.serialize<api_pb.ReferenceAndCiphertext>;
    responseDeserialize: grpc.deserialize<api_pb.ReferenceAndCiphertext>;
}
interface IEncryptionService_IDecrypt extends grpc.MethodDefinition<api_pb.ReferenceAndCiphertext, api_pb.Plaintext> {
    path: "/api.Encryption/Decrypt";
    requestStream: true;
    responseStream: true;
    requestSerialize: grpc.serialize<api_pb.ReferenceAndCiphertext>;
    requestDeserialize: grpc.deserialize<api_pb.ReferenceAndCiphertext>;
    responseSerialize: grpc.serialize<api_pb.Plaintext>;
    responseDeserialize: grpc.deserialize<api_pb.Plaintext>;
}

export const EncryptionService: IEncryptionService;

export interface IEncryptionServer {
    encrypt: grpc.handleBidiStreamingCall<api_pb.Plaintext, api_pb.ReferenceAndCiphertext>;
    decrypt: grpc.handleBidiStreamingCall<api_pb.ReferenceAndCiphertext, api_pb.Plaintext>;
}

export interface IEncryptionClient {
    encrypt(): grpc.ClientDuplexStream<api_pb.Plaintext, api_pb.ReferenceAndCiphertext>;
    encrypt(options: Partial<grpc.CallOptions>): grpc.ClientDuplexStream<api_pb.Plaintext, api_pb.ReferenceAndCiphertext>;
    encrypt(metadata: grpc.Metadata, options?: Partial<grpc.CallOptions>): grpc.ClientDuplexStream<api_pb.Plaintext, api_pb.ReferenceAndCiphertext>;
    decrypt(): grpc.ClientDuplexStream<api_pb.ReferenceAndCiphertext, api_pb.Plaintext>;
    decrypt(options: Partial<grpc.CallOptions>): grpc.ClientDuplexStream<api_pb.ReferenceAndCiphertext, api_pb.Plaintext>;
    decrypt(metadata: grpc.Metadata, options?: Partial<grpc.CallOptions>): grpc.ClientDuplexStream<api_pb.ReferenceAndCiphertext, api_pb.Plaintext>;
}

export class EncryptionClient extends grpc.Client implements IEncryptionClient {
    constructor(address: string, credentials: grpc.ChannelCredentials, options?: Partial<grpc.ClientOptions>);
    public encrypt(options?: Partial<grpc.CallOptions>): grpc.ClientDuplexStream<api_pb.Plaintext, api_pb.ReferenceAndCiphertext>;
    public encrypt(metadata?: grpc.Metadata, options?: Partial<grpc.CallOptions>): grpc.ClientDuplexStream<api_pb.Plaintext, api_pb.ReferenceAndCiphertext>;
    public decrypt(options?: Partial<grpc.CallOptions>): grpc.ClientDuplexStream<api_pb.ReferenceAndCiphertext, api_pb.Plaintext>;
    public decrypt(metadata?: grpc.Metadata, options?: Partial<grpc.CallOptions>): grpc.ClientDuplexStream<api_pb.ReferenceAndCiphertext, api_pb.Plaintext>;
}

interface IStorageService extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
    push: IStorageService_IPush;
    pull: IStorageService_IPull;
    stat: IStorageService_IStat;
    delete: IStorageService_IDelete;
}

interface IStorageService_IPush extends grpc.MethodDefinition<api_pb.Ciphertext, api_pb.Address> {
    path: "/api.Storage/Push";
    requestStream: true;
    responseStream: true;
    requestSerialize: grpc.serialize<api_pb.Ciphertext>;
    requestDeserialize: grpc.deserialize<api_pb.Ciphertext>;
    responseSerialize: grpc.serialize<api_pb.Address>;
    responseDeserialize: grpc.deserialize<api_pb.Address>;
}
interface IStorageService_IPull extends grpc.MethodDefinition<api_pb.Address, api_pb.Ciphertext> {
    path: "/api.Storage/Pull";
    requestStream: true;
    responseStream: true;
    requestSerialize: grpc.serialize<api_pb.Address>;
    requestDeserialize: grpc.deserialize<api_pb.Address>;
    responseSerialize: grpc.serialize<api_pb.Ciphertext>;
    responseDeserialize: grpc.deserialize<api_pb.Ciphertext>;
}
interface IStorageService_IStat extends grpc.MethodDefinition<api_pb.Address, stores_pb.StatInfo> {
    path: "/api.Storage/Stat";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<api_pb.Address>;
    requestDeserialize: grpc.deserialize<api_pb.Address>;
    responseSerialize: grpc.serialize<stores_pb.StatInfo>;
    responseDeserialize: grpc.deserialize<stores_pb.StatInfo>;
}
interface IStorageService_IDelete extends grpc.MethodDefinition<api_pb.Address, api_pb.Address> {
    path: "/api.Storage/Delete";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<api_pb.Address>;
    requestDeserialize: grpc.deserialize<api_pb.Address>;
    responseSerialize: grpc.serialize<api_pb.Address>;
    responseDeserialize: grpc.deserialize<api_pb.Address>;
}

export const StorageService: IStorageService;

export interface IStorageServer {
    push: grpc.handleBidiStreamingCall<api_pb.Ciphertext, api_pb.Address>;
    pull: grpc.handleBidiStreamingCall<api_pb.Address, api_pb.Ciphertext>;
    stat: grpc.handleUnaryCall<api_pb.Address, stores_pb.StatInfo>;
    delete: grpc.handleUnaryCall<api_pb.Address, api_pb.Address>;
}

export interface IStorageClient {
    push(): grpc.ClientDuplexStream<api_pb.Ciphertext, api_pb.Address>;
    push(options: Partial<grpc.CallOptions>): grpc.ClientDuplexStream<api_pb.Ciphertext, api_pb.Address>;
    push(metadata: grpc.Metadata, options?: Partial<grpc.CallOptions>): grpc.ClientDuplexStream<api_pb.Ciphertext, api_pb.Address>;
    pull(): grpc.ClientDuplexStream<api_pb.Address, api_pb.Ciphertext>;
    pull(options: Partial<grpc.CallOptions>): grpc.ClientDuplexStream<api_pb.Address, api_pb.Ciphertext>;
    pull(metadata: grpc.Metadata, options?: Partial<grpc.CallOptions>): grpc.ClientDuplexStream<api_pb.Address, api_pb.Ciphertext>;
    stat(request: api_pb.Address, callback: (error: grpc.ServiceError | null, response: stores_pb.StatInfo) => void): grpc.ClientUnaryCall;
    stat(request: api_pb.Address, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: stores_pb.StatInfo) => void): grpc.ClientUnaryCall;
    stat(request: api_pb.Address, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: stores_pb.StatInfo) => void): grpc.ClientUnaryCall;
    delete(request: api_pb.Address, callback: (error: grpc.ServiceError | null, response: api_pb.Address) => void): grpc.ClientUnaryCall;
    delete(request: api_pb.Address, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_pb.Address) => void): grpc.ClientUnaryCall;
    delete(request: api_pb.Address, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_pb.Address) => void): grpc.ClientUnaryCall;
}

export class StorageClient extends grpc.Client implements IStorageClient {
    constructor(address: string, credentials: grpc.ChannelCredentials, options?: Partial<grpc.ClientOptions>);
    public push(options?: Partial<grpc.CallOptions>): grpc.ClientDuplexStream<api_pb.Ciphertext, api_pb.Address>;
    public push(metadata?: grpc.Metadata, options?: Partial<grpc.CallOptions>): grpc.ClientDuplexStream<api_pb.Ciphertext, api_pb.Address>;
    public pull(options?: Partial<grpc.CallOptions>): grpc.ClientDuplexStream<api_pb.Address, api_pb.Ciphertext>;
    public pull(metadata?: grpc.Metadata, options?: Partial<grpc.CallOptions>): grpc.ClientDuplexStream<api_pb.Address, api_pb.Ciphertext>;
    public stat(request: api_pb.Address, callback: (error: grpc.ServiceError | null, response: stores_pb.StatInfo) => void): grpc.ClientUnaryCall;
    public stat(request: api_pb.Address, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: stores_pb.StatInfo) => void): grpc.ClientUnaryCall;
    public stat(request: api_pb.Address, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: stores_pb.StatInfo) => void): grpc.ClientUnaryCall;
    public delete(request: api_pb.Address, callback: (error: grpc.ServiceError | null, response: api_pb.Address) => void): grpc.ClientUnaryCall;
    public delete(request: api_pb.Address, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: api_pb.Address) => void): grpc.ClientUnaryCall;
    public delete(request: api_pb.Address, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: api_pb.Address) => void): grpc.ClientUnaryCall;
}

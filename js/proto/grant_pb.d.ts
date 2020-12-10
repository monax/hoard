// package: grant
// file: grant.proto

/* tslint:disable */
/* eslint-disable */

import * as jspb from "google-protobuf";
import * as github_com_gogo_protobuf_gogoproto_gogo_pb from "./github.com/gogo/protobuf/gogoproto/gogo_pb";

export class Grant extends jspb.Message { 

    hasSpec(): boolean;
    clearSpec(): void;
    getSpec(): Spec | undefined;
    setSpec(value?: Spec): Grant;

    getEncryptedreferences(): Uint8Array | string;
    getEncryptedreferences_asU8(): Uint8Array;
    getEncryptedreferences_asB64(): string;
    setEncryptedreferences(value: Uint8Array | string): Grant;

    getVersion(): number;
    setVersion(value: number): Grant;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Grant.AsObject;
    static toObject(includeInstance: boolean, msg: Grant): Grant.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Grant, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Grant;
    static deserializeBinaryFromReader(message: Grant, reader: jspb.BinaryReader): Grant;
}

export namespace Grant {
    export type AsObject = {
        spec?: Spec.AsObject,
        encryptedreferences: Uint8Array | string,
        version: number,
    }
}

export class Spec extends jspb.Message { 

    hasPlaintext(): boolean;
    clearPlaintext(): void;
    getPlaintext(): PlaintextSpec | undefined;
    setPlaintext(value?: PlaintextSpec): Spec;


    hasSymmetric(): boolean;
    clearSymmetric(): void;
    getSymmetric(): SymmetricSpec | undefined;
    setSymmetric(value?: SymmetricSpec): Spec;


    hasOpenpgp(): boolean;
    clearOpenpgp(): void;
    getOpenpgp(): OpenPGPSpec | undefined;
    setOpenpgp(value?: OpenPGPSpec): Spec;

    getLinknonce(): Uint8Array | string;
    getLinknonce_asU8(): Uint8Array;
    getLinknonce_asB64(): string;
    setLinknonce(value: Uint8Array | string): Spec;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Spec.AsObject;
    static toObject(includeInstance: boolean, msg: Spec): Spec.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Spec, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Spec;
    static deserializeBinaryFromReader(message: Spec, reader: jspb.BinaryReader): Spec;
}

export namespace Spec {
    export type AsObject = {
        plaintext?: PlaintextSpec.AsObject,
        symmetric?: SymmetricSpec.AsObject,
        openpgp?: OpenPGPSpec.AsObject,
        linknonce: Uint8Array | string,
    }
}

export class PlaintextSpec extends jspb.Message { 

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): PlaintextSpec.AsObject;
    static toObject(includeInstance: boolean, msg: PlaintextSpec): PlaintextSpec.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: PlaintextSpec, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): PlaintextSpec;
    static deserializeBinaryFromReader(message: PlaintextSpec, reader: jspb.BinaryReader): PlaintextSpec;
}

export namespace PlaintextSpec {
    export type AsObject = {
    }
}

export class SymmetricSpec extends jspb.Message { 
    getPublicid(): string;
    setPublicid(value: string): SymmetricSpec;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): SymmetricSpec.AsObject;
    static toObject(includeInstance: boolean, msg: SymmetricSpec): SymmetricSpec.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: SymmetricSpec, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): SymmetricSpec;
    static deserializeBinaryFromReader(message: SymmetricSpec, reader: jspb.BinaryReader): SymmetricSpec;
}

export namespace SymmetricSpec {
    export type AsObject = {
        publicid: string,
    }
}

export class OpenPGPSpec extends jspb.Message { 
    getPublickey(): string;
    setPublickey(value: string): OpenPGPSpec;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): OpenPGPSpec.AsObject;
    static toObject(includeInstance: boolean, msg: OpenPGPSpec): OpenPGPSpec.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: OpenPGPSpec, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): OpenPGPSpec;
    static deserializeBinaryFromReader(message: OpenPGPSpec, reader: jspb.BinaryReader): OpenPGPSpec;
}

export namespace OpenPGPSpec {
    export type AsObject = {
        publickey: string,
    }
}

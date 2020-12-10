// package: api
// file: api.proto

/* tslint:disable */
/* eslint-disable */

import * as jspb from "google-protobuf";
import * as grant_pb from "./grant_pb";
import * as reference_pb from "./reference_pb";
import * as stores_pb from "./stores_pb";

export class GrantAndGrantSpec extends jspb.Message { 

    hasGrant(): boolean;
    clearGrant(): void;
    getGrant(): grant_pb.Grant | undefined;
    setGrant(value?: grant_pb.Grant): GrantAndGrantSpec;


    hasGrantspec(): boolean;
    clearGrantspec(): void;
    getGrantspec(): grant_pb.Spec | undefined;
    setGrantspec(value?: grant_pb.Spec): GrantAndGrantSpec;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): GrantAndGrantSpec.AsObject;
    static toObject(includeInstance: boolean, msg: GrantAndGrantSpec): GrantAndGrantSpec.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: GrantAndGrantSpec, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): GrantAndGrantSpec;
    static deserializeBinaryFromReader(message: GrantAndGrantSpec, reader: jspb.BinaryReader): GrantAndGrantSpec;
}

export namespace GrantAndGrantSpec {
    export type AsObject = {
        grant?: grant_pb.Grant.AsObject,
        grantspec?: grant_pb.Spec.AsObject,
    }
}

export class PlaintextAndGrantSpec extends jspb.Message { 

    hasPlaintext(): boolean;
    clearPlaintext(): void;
    getPlaintext(): Plaintext | undefined;
    setPlaintext(value?: Plaintext): PlaintextAndGrantSpec;


    hasGrantspec(): boolean;
    clearGrantspec(): void;
    getGrantspec(): grant_pb.Spec | undefined;
    setGrantspec(value?: grant_pb.Spec): PlaintextAndGrantSpec;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): PlaintextAndGrantSpec.AsObject;
    static toObject(includeInstance: boolean, msg: PlaintextAndGrantSpec): PlaintextAndGrantSpec.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: PlaintextAndGrantSpec, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): PlaintextAndGrantSpec;
    static deserializeBinaryFromReader(message: PlaintextAndGrantSpec, reader: jspb.BinaryReader): PlaintextAndGrantSpec;
}

export namespace PlaintextAndGrantSpec {
    export type AsObject = {
        plaintext?: Plaintext.AsObject,
        grantspec?: grant_pb.Spec.AsObject,
    }
}

export class ReferenceAndGrantSpec extends jspb.Message { 

    hasReference(): boolean;
    clearReference(): void;
    getReference(): reference_pb.Ref | undefined;
    setReference(value?: reference_pb.Ref): ReferenceAndGrantSpec;


    hasGrantspec(): boolean;
    clearGrantspec(): void;
    getGrantspec(): grant_pb.Spec | undefined;
    setGrantspec(value?: grant_pb.Spec): ReferenceAndGrantSpec;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ReferenceAndGrantSpec.AsObject;
    static toObject(includeInstance: boolean, msg: ReferenceAndGrantSpec): ReferenceAndGrantSpec.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ReferenceAndGrantSpec, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ReferenceAndGrantSpec;
    static deserializeBinaryFromReader(message: ReferenceAndGrantSpec, reader: jspb.BinaryReader): ReferenceAndGrantSpec;
}

export namespace ReferenceAndGrantSpec {
    export type AsObject = {
        reference?: reference_pb.Ref.AsObject,
        grantspec?: grant_pb.Spec.AsObject,
    }
}

export class Header extends jspb.Message { 
    getSalt(): Uint8Array | string;
    getSalt_asU8(): Uint8Array;
    getSalt_asB64(): string;
    setSalt(value: Uint8Array | string): Header;

    getData(): Uint8Array | string;
    getData_asU8(): Uint8Array;
    getData_asB64(): string;
    setData(value: Uint8Array | string): Header;

    getChunksize(): number;
    setChunksize(value: number): Header;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Header.AsObject;
    static toObject(includeInstance: boolean, msg: Header): Header.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Header, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Header;
    static deserializeBinaryFromReader(message: Header, reader: jspb.BinaryReader): Header;
}

export namespace Header {
    export type AsObject = {
        salt: Uint8Array | string,
        data: Uint8Array | string,
        chunksize: number,
    }
}

export class Plaintext extends jspb.Message { 
    getBody(): Uint8Array | string;
    getBody_asU8(): Uint8Array;
    getBody_asB64(): string;
    setBody(value: Uint8Array | string): Plaintext;


    hasHead(): boolean;
    clearHead(): void;
    getHead(): Header | undefined;
    setHead(value?: Header): Plaintext;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Plaintext.AsObject;
    static toObject(includeInstance: boolean, msg: Plaintext): Plaintext.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Plaintext, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Plaintext;
    static deserializeBinaryFromReader(message: Plaintext, reader: jspb.BinaryReader): Plaintext;
}

export namespace Plaintext {
    export type AsObject = {
        body: Uint8Array | string,
        head?: Header.AsObject,
    }
}

export class Ciphertext extends jspb.Message { 
    getEncrypteddata(): Uint8Array | string;
    getEncrypteddata_asU8(): Uint8Array;
    getEncrypteddata_asB64(): string;
    setEncrypteddata(value: Uint8Array | string): Ciphertext;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Ciphertext.AsObject;
    static toObject(includeInstance: boolean, msg: Ciphertext): Ciphertext.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Ciphertext, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Ciphertext;
    static deserializeBinaryFromReader(message: Ciphertext, reader: jspb.BinaryReader): Ciphertext;
}

export namespace Ciphertext {
    export type AsObject = {
        encrypteddata: Uint8Array | string,
    }
}

export class ReferenceAndCiphertext extends jspb.Message { 

    hasReference(): boolean;
    clearReference(): void;
    getReference(): reference_pb.Ref | undefined;
    setReference(value?: reference_pb.Ref): ReferenceAndCiphertext;


    hasCiphertext(): boolean;
    clearCiphertext(): void;
    getCiphertext(): Ciphertext | undefined;
    setCiphertext(value?: Ciphertext): ReferenceAndCiphertext;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ReferenceAndCiphertext.AsObject;
    static toObject(includeInstance: boolean, msg: ReferenceAndCiphertext): ReferenceAndCiphertext.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ReferenceAndCiphertext, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ReferenceAndCiphertext;
    static deserializeBinaryFromReader(message: ReferenceAndCiphertext, reader: jspb.BinaryReader): ReferenceAndCiphertext;
}

export namespace ReferenceAndCiphertext {
    export type AsObject = {
        reference?: reference_pb.Ref.AsObject,
        ciphertext?: Ciphertext.AsObject,
    }
}

export class Address extends jspb.Message { 
    getAddress(): Uint8Array | string;
    getAddress_asU8(): Uint8Array;
    getAddress_asB64(): string;
    setAddress(value: Uint8Array | string): Address;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Address.AsObject;
    static toObject(includeInstance: boolean, msg: Address): Address.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Address, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Address;
    static deserializeBinaryFromReader(message: Address, reader: jspb.BinaryReader): Address;
}

export namespace Address {
    export type AsObject = {
        address: Uint8Array | string,
    }
}

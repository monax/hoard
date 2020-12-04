// package: reference
// file: reference.proto

/* tslint:disable */
/* eslint-disable */

import * as jspb from "google-protobuf";

export class Ref extends jspb.Message { 
    getAddress(): Uint8Array | string;
    getAddress_asU8(): Uint8Array;
    getAddress_asB64(): string;
    setAddress(value: Uint8Array | string): Ref;

    getSecretkey(): Uint8Array | string;
    getSecretkey_asU8(): Uint8Array;
    getSecretkey_asB64(): string;
    setSecretkey(value: Uint8Array | string): Ref;

    getSalt(): Uint8Array | string;
    getSalt_asU8(): Uint8Array;
    getSalt_asB64(): string;
    setSalt(value: Uint8Array | string): Ref;

    getType(): Ref.RefType;
    setType(value: Ref.RefType): Ref;

    getSize(): number;
    setSize(value: number): Ref;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Ref.AsObject;
    static toObject(includeInstance: boolean, msg: Ref): Ref.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Ref, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Ref;
    static deserializeBinaryFromReader(message: Ref, reader: jspb.BinaryReader): Ref;
}

export namespace Ref {
    export type AsObject = {
        address: Uint8Array | string,
        secretkey: Uint8Array | string,
        salt: Uint8Array | string,
        type: Ref.RefType,
        size: number,
    }

    export enum RefType {
    BODY = 0,
    HEADER = 1,
    LINK = 2,
    }

}

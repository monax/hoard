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
    getVersion(): number;
    setVersion(value: number): Ref;
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
        version: number,
        type: Ref.RefType,
        size: number,
    }

    export enum RefType {
    BODY = 0,
    HEADER = 1,
    LINK = 2,
    }

}

export class RefsWithNonce extends jspb.Message { 
    clearRefsList(): void;
    getRefsList(): Array<Ref>;
    setRefsList(value: Array<Ref>): RefsWithNonce;
    addRefs(value?: Ref, index?: number): Ref;
    getNonce(): Uint8Array | string;
    getNonce_asU8(): Uint8Array;
    getNonce_asB64(): string;
    setNonce(value: Uint8Array | string): RefsWithNonce;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): RefsWithNonce.AsObject;
    static toObject(includeInstance: boolean, msg: RefsWithNonce): RefsWithNonce.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: RefsWithNonce, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): RefsWithNonce;
    static deserializeBinaryFromReader(message: RefsWithNonce, reader: jspb.BinaryReader): RefsWithNonce;
}

export namespace RefsWithNonce {
    export type AsObject = {
        refsList: Array<Ref.AsObject>,
        nonce: Uint8Array | string,
    }
}

export class Link extends jspb.Message { 

    hasHeader(): boolean;
    clearHeader(): void;
    getHeader(): Ref | undefined;
    setHeader(value?: Ref): Link;
    clearBodyList(): void;
    getBodyList(): Array<Ref>;
    setBodyList(value: Array<Ref>): Link;
    addBody(value?: Ref, index?: number): Ref;

    hasTrailer(): boolean;
    clearTrailer(): void;
    getTrailer(): Ref | undefined;
    setTrailer(value?: Ref): Link;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Link.AsObject;
    static toObject(includeInstance: boolean, msg: Link): Link.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Link, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Link;
    static deserializeBinaryFromReader(message: Link, reader: jspb.BinaryReader): Link;
}

export namespace Link {
    export type AsObject = {
        header?: Ref.AsObject,
        bodyList: Array<Ref.AsObject>,
        trailer?: Ref.AsObject,
    }
}

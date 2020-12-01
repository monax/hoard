// package: stores
// file: stores.proto

/* tslint:disable */
/* eslint-disable */

import * as jspb from "google-protobuf";

export class StatInfo extends jspb.Message { 
    getAddress(): Uint8Array | string;
    getAddress_asU8(): Uint8Array;
    getAddress_asB64(): string;
    setAddress(value: Uint8Array | string): StatInfo;

    getExists(): boolean;
    setExists(value: boolean): StatInfo;

    getSize(): number;
    setSize(value: number): StatInfo;

    getLocation(): string;
    setLocation(value: string): StatInfo;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): StatInfo.AsObject;
    static toObject(includeInstance: boolean, msg: StatInfo): StatInfo.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: StatInfo, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): StatInfo;
    static deserializeBinaryFromReader(message: StatInfo, reader: jspb.BinaryReader): StatInfo;
}

export namespace StatInfo {
    export type AsObject = {
        address: Uint8Array | string,
        exists: boolean,
        size: number,
        location: string,
    }
}

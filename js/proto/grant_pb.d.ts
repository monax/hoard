// package: grant
// file: grant.proto

import * as jspb from "google-protobuf";
import * as github_com_gogo_protobuf_gogoproto_gogo_pb from "./github.com/gogo/protobuf/gogoproto/gogo_pb";

export class Grant extends jspb.Message {
  hasSpec(): boolean;
  clearSpec(): void;
  getSpec(): Spec | undefined;
  setSpec(value?: Spec): void;

  getEncryptedreferences(): Uint8Array | string;
  getEncryptedreferences_asU8(): Uint8Array;
  getEncryptedreferences_asB64(): string;
  setEncryptedreferences(value: Uint8Array | string): void;

  getVersion(): number;
  setVersion(value: number): void;

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
  setPlaintext(value?: PlaintextSpec): void;

  hasSymmetric(): boolean;
  clearSymmetric(): void;
  getSymmetric(): SymmetricSpec | undefined;
  setSymmetric(value?: SymmetricSpec): void;

  hasOpenpgp(): boolean;
  clearOpenpgp(): void;
  getOpenpgp(): OpenPGPSpec | undefined;
  setOpenpgp(value?: OpenPGPSpec): void;

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
  setPublicid(value: string): void;

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
  setPublickey(value: string): void;

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


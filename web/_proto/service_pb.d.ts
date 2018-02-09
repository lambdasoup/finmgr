// package: pb
// file: service.proto

import * as jspb from "google-protobuf";

export class Hello extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Hello.AsObject;
  static toObject(includeInstance: boolean, msg: Hello): Hello.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Hello, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Hello;
  static deserializeBinaryFromReader(message: Hello, reader: jspb.BinaryReader): Hello;
}

export namespace Hello {
  export type AsObject = {
    name: string,
  }
}

export class Bye extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Bye.AsObject;
  static toObject(includeInstance: boolean, msg: Bye): Bye.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Bye, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Bye;
  static deserializeBinaryFromReader(message: Bye, reader: jspb.BinaryReader): Bye;
}

export namespace Bye {
  export type AsObject = {
    name: string,
  }
}

export class Empty extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Empty.AsObject;
  static toObject(includeInstance: boolean, msg: Empty): Empty.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Empty, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Empty;
  static deserializeBinaryFromReader(message: Empty, reader: jspb.BinaryReader): Empty;
}

export namespace Empty {
  export type AsObject = {
  }
}


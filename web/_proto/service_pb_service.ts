// package: pb
// file: service.proto

import * as service_pb from "./service_pb";
export class Service {
  static serviceName = "pb.Service";
}
export namespace Service {
  export class SayHello {
    static readonly methodName = "SayHello";
    static readonly service = Service;
    static readonly requestStream = false;
    static readonly responseStream = false;
    static readonly requestType = service_pb.Hello;
    static readonly responseType = service_pb.Bye;
  }
}

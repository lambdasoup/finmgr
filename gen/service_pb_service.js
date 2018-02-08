// package: pb
// file: service.proto

var jspb = require("google-protobuf");
var service_pb = require("./service_pb");
var Service = {
  serviceName: "pb.Service"
};
Service.SayHello = {
  methodName: "SayHello",
  service: Service,
  requestStream: false,
  responseStream: false,
  requestType: service_pb.Hello,
  responseType: service_pb.Bye
};
module.exports = {
  Service: Service,
};


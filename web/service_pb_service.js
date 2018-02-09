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
Service.GetHellos = {
  methodName: "GetHellos",
  service: Service,
  requestStream: false,
  responseStream: true,
  requestType: service_pb.Empty,
  responseType: service_pb.Hello
};
module.exports = {
  Service: Service,
};


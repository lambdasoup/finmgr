// package: pb
// file: example.proto

var jspb = require("google-protobuf");
var example_pb = require("./example_pb");
var Greeter = {
  serviceName: "pb.Greeter"
};
Greeter.SayHello = {
  methodName: "SayHello",
  service: Greeter,
  requestStream: false,
  responseStream: false,
  requestType: example_pb.HelloRequest,
  responseType: example_pb.HelloReply
};
module.exports = {
  Greeter: Greeter,
};


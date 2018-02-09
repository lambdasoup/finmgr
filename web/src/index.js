var app = Elm.Main.fullscreen()

var service = require("../service_pb_service")
var pb = require("../service_pb")
var grpc = require("grpc-web-client")

var host = "http://localhost:8080";

function sendHello() {
  var hello = new pb.Hello();
  hello.setName("Alice");
  grpc.grpc.unary(service.Service.SayHello, {
    request: hello,
    host: host,
    onEnd: function(res) {
      console.log("res", res);
    }
  });
}

function getHellos() {
  var empty = new pb.Empty();
  const client = grpc.grpc.client(service.Service.GetHellos, {
    host: host,
  });
  client.onHeaders(function(headers)  {
    console.log("query.onHeaders", headers);
  });
  client.onMessage(function(message) {
    console.log("query.onMessage", message);
  });
  client.onEnd(function(code, msg, trailers) {
    console.log("query.onEnd", code, msg, trailers);
  });
  client.start();
  client.send(empty);
}

sendHello()
getHellos()

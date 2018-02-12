var app = Elm.Main.fullscreen()

var service = require("../service_pb_service")
var jspb = require("google-protobuf")
var pb = require("../service_pb")
var grpc = require("grpc-web-client")

var host = window.location.protocol + "//" + window.location.host;

app.ports.hello.subscribe(function(msg) {
  var hello = new pb.Hello();
  hello.setName(msg.name);
  grpc.grpc.unary(service.Service.SayHello, {
    request: hello,
    host: host,
    onEnd: function(res) {
      app.ports.reply.send(res.message.toObject());
    }
  });
});

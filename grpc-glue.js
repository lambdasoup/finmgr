var app = Elm.Main.fullscreen()
var service = require("./service_pb_service")
var pb = require("./service_pb")
var jspb = require("google-protobuf")
var grpc = require("grpc-web-client")

var host = window.location.protocol + "//" + window.location.host;

app.ports.getAccounts.subscribe(function(msg) {
  var accountInfo = new pb.AccountInfo();
  accountInfo.setId(msg.id);
  accountInfo.setPin(msg.pin);
  accountInfo.setBlz(msg.blz);
  console.log(msg)
  grpc.grpc.unary(service.Service.GetAccounts, {
    request: accountInfo,
    host: host,
    onEnd: function(res) {
      app.ports.reply.send(res.message.toObject());
    }
  });
});

var service = require("./service_pb_service")
var pb = require("./service_pb")
var jspb = require("google-protobuf")
var grpc = require("grpc-web-client")

export function connect(app) {
  var host = window.location.protocol + "//" + window.location.host;
  app.ports.addBank.subscribe(function(msg) {
    var bank = new pb.Bank();
    bank.setId(msg.id);
    bank.setPin(msg.pin);
    bank.setBlz(msg.blz);
    grpc.grpc.unary(service.AccountService.AddBank, {
      request: bank,
      host: host,
      onEnd: function(res) {
        // TODO log error?;
      }
    });
  });

  app.ports.getUserEmpty.subscribe(function(msg) {
    var empty = new pb.Empty();
    grpc.grpc.unary(service.UserService.GetUser, {
      request: empty,
      host: host,
      onEnd: function(res) {
        if (res.status == 2) {
          // TODO error
          console.log(res.statusMessage);
        } else {
          app.ports.setUser.send(res.message.toObject());
        }
      }
    });
  });
}

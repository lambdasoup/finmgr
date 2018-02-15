//# sourceURL=app.js
'use strict';

var elm = require("../elm-src/Main.elm")
var app = elm.Main.fullscreen()
var service = require("./service_pb_service")
var pb = require("./service_pb")
var jspb = require("google-protobuf")
var grpc = require("grpc-web-client")
var push = require("./push")

var host = window.location.protocol + "//" + window.location.host;

app.ports.getAccounts.subscribe(function(msg) {
  var accountInfo = new pb.AccountInfo();
  accountInfo.setId(msg.id);
  accountInfo.setPin(msg.pin);
  accountInfo.setBlz(msg.blz);
  grpc.grpc.unary(service.Service.GetAccounts, {
    request: accountInfo,
    host: host,
    onEnd: function(res) {
      app.ports.reply.send(res.message.toObject());
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

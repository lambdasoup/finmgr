import * as Elm from './Main'

import {grpc} from "grpc-web-client";
import {Service} from "../_proto/service_pb_service";
import {Hello, Bye} from "../_proto/service_pb";

let app = Elm.Main.fullscreen()
app.ports.hello.subscribe(name => console.log(`Hello ${name}!!`))
app.ports.reply.send(12345)

const host = "http://localhost:8080";

function sendHello() {
  const hello = new Hello();
  hello.setName("Alice");
  grpc.unary(Service.SayHello, {
    request: hello,
    host: host,
    onEnd: res => {
      const { status, statusMessage, headers, message, trailers } = res;
      console.log("onEnd.status", status, statusMessage);
      console.log("onEnd.headers", headers);
      if (status === grpc.Code.OK && message) {
        console.log("onEnd.message", message.toObject());
      }
      console.log("onEnd.trailers", trailers);
    }
  });
}

sendHello();


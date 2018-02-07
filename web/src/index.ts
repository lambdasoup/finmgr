import * as Elm from './Main'

import {grpc} from "grpc-web-client";
import {BookService} from "../_proto/proto/examplecom/library/book_service_pb_service";
import {QueryBooksRequest, Book, GetBookRequest} from "../_proto/proto/examplecom/library/book_service_pb";

let app = Elm.Main.fullscreen()
app.ports.hello.subscribe(name => console.log(`Hello ${name}!!`))
app.ports.reply.send(12345)

const host = "http://localhost:8080/api";

function getBook() {
  const getBookRequest = new GetBookRequest();
  getBookRequest.setIsbn(60929871);
  grpc.unary(BookService.GetBook, {
    request: getBookRequest,
    host: host,
    onEnd: res => {
      const { status, statusMessage, headers, message, trailers } = res;
      console.log("getBook.onEnd.status", status, statusMessage);
      console.log("getBook.onEnd.headers", headers);
      if (status === grpc.Code.OK && message) {
        console.log("getBook.onEnd.message", message.toObject());
      }
      console.log("getBook.onEnd.trailers", trailers);
      queryBooks();
    }
  });
}

getBook();

function queryBooks() {
  const queryBooksRequest = new QueryBooksRequest();
  queryBooksRequest.setAuthorPrefix("Geor");
  const client = grpc.client(BookService.QueryBooks, {
    host: host,
  });
  client.onHeaders((headers: grpc.Metadata) => {
     console.log("queryBooks.onHeaders", headers);
  });
  //  client.onMessage((message: Book) => {
  //    console.log("queryBooks.onMessage", message.toObject());
  //  });
  client.onEnd((code: grpc.Code, msg: string, trailers: grpc.Metadata) => {
    console.log("queryBooks.onEnd", code, msg, trailers);
  });
  client.start();
  client.send(queryBooksRequest);
}


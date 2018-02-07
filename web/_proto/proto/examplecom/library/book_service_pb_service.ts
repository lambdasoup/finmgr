// package: examplecom.library
// file: proto/examplecom/library/book_service.proto

import * as proto_examplecom_library_book_service_pb from "../../../proto/examplecom/library/book_service_pb";
export class BookService {
  static serviceName = "examplecom.library.BookService";
}
export namespace BookService {
  export class GetBook {
    static readonly methodName = "GetBook";
    static readonly service = BookService;
    static readonly requestStream = false;
    static readonly responseStream = false;
    static readonly requestType = proto_examplecom_library_book_service_pb.GetBookRequest;
    static readonly responseType = proto_examplecom_library_book_service_pb.Book;
  }
  export class QueryBooks {
    static readonly methodName = "QueryBooks";
    static readonly service = BookService;
    static readonly requestStream = false;
    static readonly responseStream = true;
    static readonly requestType = proto_examplecom_library_book_service_pb.QueryBooksRequest;
    static readonly responseType = proto_examplecom_library_book_service_pb.Book;
  }
}

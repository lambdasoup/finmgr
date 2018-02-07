rm ../gen/*

# go
protoc --go_out=plugins=grpc:../gen/ *.proto

# elm (types)
protoc --elm_out=plugins=grpc:../gen/ *.proto

# js (services & types)
protoc \
  --plugin=protoc-gen-js_service=./node_modules/.bin/protoc-gen-js_service \
  --js_out=import_style=commonjs,binary:../gen/ \
  --js_service_out=../gen/ \
  -I . \
  *.proto

# ts
protoc \
  --plugin=protoc-gen-ts=./node_modules/.bin/protoc-gen-ts \
  --js_out=import_style=commonjs,binary:./_proto \
  --ts_out=service=true:./_proto \
  ./proto/examplecom/library/book_service.proto

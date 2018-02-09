rm ../gen/*
rm _proto/*

# go
protoc --go_out=plugins=grpc:../gen/ service.proto

# elm (types)
protoc --elm_out=plugins=grpc:../gen/ service.proto

# js (services & types)
protoc \
  --plugin=protoc-gen-js_service=./node_modules/.bin/protoc-gen-js_service \
  --js_out=import_style=commonjs,binary:. \
  --js_service_out=. \
  -I . \
  service.proto

# ts
protoc \
  --plugin=protoc-gen-ts=./node_modules/.bin/protoc-gen-ts \
  --js_out=import_style=commonjs,binary:./_proto \
  --ts_out=service=true:./_proto \
  ./service.proto

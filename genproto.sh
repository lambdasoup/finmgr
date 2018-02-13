protoc \
  --go_out=plugins=grpc:. \
  --elm_out=plugins=grpc:elm-src \
  --plugin=protoc-gen-js_service=./node_modules/.bin/protoc-gen-js_service \
  --js_out=import_style=commonjs,binary:. \
  --js_service_out=. \
  service.proto \
  2>&1

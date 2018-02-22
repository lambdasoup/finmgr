
proto = js-src/service_pb.js js-src/service_pb_service.js service.pb.go elm-src/Service.elm
elm = js-src/elm.js
js = app/index.js app/index.js.map

# all: $(js)
all: $(js)

clean:
	rm -f $(proto)
	rm -f $(elm)
	rm -f $(js)

$(js): $(elm) $(proto) js-src/*.js webpack.config.js
	webpack

$(elm): $(proto) elm-src/*.elm
	elm-make elm-src/Main.elm --output $(elm)

$(proto): service.proto
	protoc \
	--go_out=plugins=grpc:. \
  --elm_out=plugins=grpc:elm-src \
  --plugin=protoc-gen-js_service=./node_modules/.bin/protoc-gen-js_service \
  --js_out=import_style=commonjs,binary:js-src \
  --js_service_out=js-src \
  service.proto \

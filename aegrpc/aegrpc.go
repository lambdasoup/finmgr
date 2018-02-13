package aegrpc

import (
	"context"
	"net/http"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/appengine"
	"google.golang.org/grpc"
)

var contextKey = "appengine context"

func HandlePath(p string, s *grpc.Server) {
	wrappedGrpc := grpcweb.WrapServer(s)
	http.HandleFunc(p, (func(resp http.ResponseWriter, req *http.Request) {
		ctx := appengine.NewContext(req)

		wrappedreq := req.WithContext(context.WithValue(req.Context(), &contextKey, ctx))

		wrappedGrpc.ServeHTTP(resp, wrappedreq)
	}))
}

func NewAppengineContext(ctx context.Context) context.Context {
	return ctx.Value(&contextKey).(context.Context)
}

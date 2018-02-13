package aegrpc

import (
	"context"
	"net/http"

	"google.golang.org/appengine"
)

var contextKey = "appengine context"

// NewAppengineHandlerFunc wraps the given http.HandlerFunc for AppEngine
// requests
func NewAppengineHandlerFunc(hf http.HandlerFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		ctx := appengine.NewContext(req)
		wreq := req.WithContext(context.WithValue(req.Context(), &contextKey, ctx))
		hf.ServeHTTP(resp, wreq)
	}
}

// NewAppengineContext returns the (native) AppEngine request context for a
// parent context.Context
func NewAppengineContext(ctx context.Context) context.Context {
	return ctx.Value(&contextKey).(context.Context)
}

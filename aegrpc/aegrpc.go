package aegrpc

import (
	"context"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/user"
)

var contextKey = "appengine context"
var userKey = "appengine user"

// NewAppengineHandlerFunc wraps the given http.HandlerFunc for AppEngine
// requests
func NewAppengineHandlerFunc(hf http.HandlerFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		ctx := appengine.NewContext(req)
		wctx := context.WithValue(req.Context(), &contextKey, ctx)

		gu := user.Current(ctx)
		uk := MakeUserKey(ctx, gu.ID)
		wctx = context.WithValue(wctx, &userKey, uk)

		wreq := req.WithContext(wctx)
		hf.ServeHTTP(resp, wreq)
	}
}

// NewAppengineContext returns the (native) AppEngine request context for a
// parent context.Context
func NewAppengineContext(ctx context.Context) context.Context {
	return ctx.Value(&contextKey).(context.Context)
}

// GetUserKey returns the datastore user key
func GetUserKey(ctx context.Context) *datastore.Key {
	return ctx.Value(&userKey).(*datastore.Key)
}

func MakeUserKey(ctx context.Context, id string) *datastore.Key {
	return datastore.NewKey(ctx, "User", id, 0, nil)
}

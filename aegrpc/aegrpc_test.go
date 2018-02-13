package aegrpc_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lambdasoup/finmgr/aegrpc"
	"google.golang.org/appengine/aetest"
)

type testHandler struct{}

func TestAppengineContextWrapping(t *testing.T) {
	aeinst, err := aetest.NewInstance(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer aeinst.Close()

	req, _ := aeinst.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()

	hf := func(resp http.ResponseWriter, req *http.Request) {
		actx := aegrpc.NewAppengineContext(req.Context())
		if actx == nil {
			t.Fail()
		}
	}

	aegrpc.NewAppengineHandlerFunc(hf).ServeHTTP(resp, req)
}

package finmgr

import (
	"net/http"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type server struct{}

func (s *server) SayHello(ctx context.Context, in *Hello) (*Bye, error) {
	grpc.SendHeader(ctx, metadata.Pairs("Pre-Response-Metadata", "Is-sent-as-headers-unary"))
	grpc.SetTrailer(ctx, metadata.Pairs("Post-Response-Metadata", "Is-sent-as-trailers-unary"))

	return &Bye{Name: "Bye " + in.Name}, nil
}

func init() {
	s := grpc.NewServer()
	RegisterServiceServer(s, &server{})

	wrappedGrpc := grpcweb.WrapServer(s)
	http.HandleFunc("/finmgr.Service/", (func(resp http.ResponseWriter, req *http.Request) {
		ctx := appengine.NewContext(req)
		log.Debugf(ctx, "got req")
		wrappedGrpc.ServeHTTP(resp, req)
	}))
}

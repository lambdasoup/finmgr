package web

import (
	pb "github.com/lambdasoup/finmgr/gen"
	"net/http"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello mplements elloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func init() {
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})

	wrappedGrpc := grpcweb.WrapServer(s)
	http.HandleFunc("/api/", (func(resp http.ResponseWriter, req *http.Request) {
		if wrappedGrpc.IsGrpcWebRequest(req) {
			wrappedGrpc.ServeHTTP(resp, req)
		}
		// Fall back to other servers.
		http.DefaultServeMux.ServeHTTP(resp, req)
	}))
}

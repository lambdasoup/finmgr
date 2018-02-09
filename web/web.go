package web

import (
	"net/http"

	pb "github.com/lambdasoup/finmgr/gen"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type server struct{}

func (s *server) SayHello(ctx context.Context, in *pb.Hello) (*pb.Bye, error) {
	grpc.SendHeader(ctx, metadata.Pairs("Pre-Response-Metadata", "Is-sent-as-headers-unary"))
	grpc.SetTrailer(ctx, metadata.Pairs("Post-Response-Metadata", "Is-sent-as-trailers-unary"))

	return &pb.Bye{Name: "Bye " + in.Name}, nil
}

func (s *server) GetHellos(in *pb.Empty, stream pb.Service_GetHellosServer) error {
	stream.SendHeader(metadata.Pairs("Pre-Response-Metadata", "Is-sent-as-headers-stream"))

	stream.Send(&pb.Hello{"hi"})
	stream.Send(&pb.Hello{"hola"})
	stream.Send(&pb.Hello{"hello"})
	stream.Send(&pb.Hello{"hey"})
	stream.Send(&pb.Hello{"hallo"})

	stream.SetTrailer(metadata.Pairs("Post-Response-Metadata", "Is-sent-as-trailers-stream"))
	return nil
}

func init() {
	s := grpc.NewServer()
	pb.RegisterServiceServer(s, &server{})

	wrappedGrpc := grpcweb.WrapServer(s)
	http.HandleFunc("/pb.Service/", (func(resp http.ResponseWriter, req *http.Request) {
		ctx := appengine.NewContext(req)
		log.Debugf(ctx, "got req")
		wrappedGrpc.ServeHTTP(resp, req)
	}))
}

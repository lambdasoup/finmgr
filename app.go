package finmgr

import (
	"fmt"
	"net/http"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/mitch000001/go-hbci/client"
	"github.com/mitch000001/go-hbci/domain"
	https "github.com/mitch000001/go-hbci/transport/https"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type server struct{}

var contextKey = "appengine context"

func (s *server) GetAccounts(ctx context.Context, in *AccountInfo) (*Accounts, error) {
	grpc.SendHeader(ctx, metadata.Pairs("Pre-Response-Metadata", "Is-sent-as-headers-unary"))
	grpc.SetTrailer(ctx, metadata.Pairs("Post-Response-Metadata", "Is-sent-as-trailers-unary"))

	actx := ctx.Value(&contextKey).(context.Context)
	log.Debugf(actx, "in %v", in)

	config := client.Config{
		AccountID:   in.GetId(),
		BankID:      in.GetBlz(),
		PIN:         in.GetPin(),
		HBCIVersion: domain.FINTSVersion300,
		Transport:   https.NewNonDefault(urlfetch.Client(actx)),
	}

	c, err := client.New(config)
	if err != nil {
		panic(err)
	}

	accounts, err := c.Accounts()

	return &Accounts{Info: fmt.Sprintf("TODO accounts info %v", accounts)}, err
}

func init() {
	s := grpc.NewServer()
	RegisterServiceServer(s, &server{})

	wrappedGrpc := grpcweb.WrapServer(s)
	http.HandleFunc("/finmgr.Service/", (func(resp http.ResponseWriter, req *http.Request) {
		ctx := appengine.NewContext(req)

		wrappedreq := req.WithContext(context.WithValue(req.Context(), &contextKey, ctx))

		wrappedGrpc.ServeHTTP(resp, wrappedreq)
	}))
}

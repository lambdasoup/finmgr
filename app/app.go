package app

import (
	"net/http"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/lambdasoup/finmgr"
	"github.com/lambdasoup/finmgr/account"
	"github.com/lambdasoup/finmgr/aegrpc"
	"github.com/lambdasoup/finmgr/user"
	"github.com/lambdasoup/finmgr/webpush"
	"google.golang.org/grpc"
)

func init() {
	s := grpc.NewServer()

	finmgr.RegisterUserServiceServer(s, user.NewServer())
	finmgr.RegisterAccountServiceServer(s, account.NewServer())

	webGrpc := grpcweb.WrapServer(s)
	http.HandleFunc("/finmgr.UserService/", aegrpc.NewAppengineHandlerFunc(webGrpc.ServeHTTP))
	http.HandleFunc("/finmgr.AccountService/", aegrpc.NewAppengineHandlerFunc(webGrpc.ServeHTTP))

	http.HandleFunc("/web-push/publicKey", webpush.GetPublicKey)

	http.HandleFunc("/worker", account.UpdateAccounts)
}

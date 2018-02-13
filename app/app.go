package app

import (
	"github.com/lambdasoup/finmgr"
	"github.com/lambdasoup/finmgr/account"
	"github.com/lambdasoup/finmgr/aegrpc"
	"github.com/lambdasoup/finmgr/user"
	"google.golang.org/grpc"
)

func init() {
	s := grpc.NewServer()

	finmgr.RegisterUserServiceServer(s, user.NewServer())
	finmgr.RegisterAccountServiceServer(s, account.NewServer())

	aegrpc.HandlePath("/finmgr.UserService/", s)
	aegrpc.HandlePath("/finmgr.AccountService/", s)
}

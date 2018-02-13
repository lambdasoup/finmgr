package user

import (
	"github.com/lambdasoup/finmgr"
	"github.com/lambdasoup/finmgr/aegrpc"
	"golang.org/x/net/context"
	"google.golang.org/appengine/user"
)

func (s *server) GetUser(ctx context.Context, in *finmgr.Empty) (*finmgr.User, error) {
	actx := aegrpc.NewAppengineContext(ctx)

	u := user.Current(actx)
	lu, err := user.LogoutURL(actx, "/")

	return &finmgr.User{Email: u.Email, LogoutUrl: lu}, err
}

type server struct{}

// NewServer returns a new implementation for a UserServiceServer
func NewServer() finmgr.UserServiceServer {
	return new(server)
}

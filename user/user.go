package user

import (
	"github.com/lambdasoup/finmgr"
	"github.com/lambdasoup/finmgr/aegrpc"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	aeuser "google.golang.org/appengine/user"
)

type server struct{}

type user struct {
}

func (s *server) GetUser(ctx context.Context, in *finmgr.Empty) (*finmgr.User, error) {
	actx := aegrpc.NewAppengineContext(ctx)

	lu, err := aeuser.LogoutURL(actx, "/")

	gu := aeuser.Current(actx)
	return &finmgr.User{Email: gu.Email, LogoutUrl: lu}, err
}

func getUser(ctx context.Context) (*user, error) {
	uk := aegrpc.GetUserKey(ctx)

	// get user from db, create if not exist
	u := user{}
	err := datastore.Get(ctx, uk, &u)
	if err == datastore.ErrNoSuchEntity {
		_, err = datastore.Put(ctx, uk, &u)
	}

	return &u, err
}

// NewServer returns a new implementation for a UserServiceServer
func NewServer() finmgr.UserServiceServer {
	return new(server)
}

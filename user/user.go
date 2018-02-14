package user

import (
	"github.com/lambdasoup/finmgr"
	"github.com/lambdasoup/finmgr/aegrpc"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	aeuser "google.golang.org/appengine/user"
)

type server struct{}

type user struct {
}

func (s *server) GetUser(ctx context.Context, in *finmgr.Empty) (*finmgr.User, error) {
	actx := aegrpc.NewAppengineContext(ctx)

	// use google user id as key
	gu := aeuser.Current(actx)
	uk := datastore.NewKey(actx, "User", gu.ID, 0, nil)

	log.Debugf(actx, "uk: %v", uk)

	// get user from db, create if not exist
	u := user{}
	err := datastore.Get(actx, uk, &u)
	if err == datastore.ErrNoSuchEntity {
		log.Debugf(actx, "trying put")
		_, err = datastore.Put(actx, uk, &u)
	}

	// bail on get/put err
	if err != nil {
		return nil, err
	}

	lu, err := aeuser.LogoutURL(actx, "/")

	return &finmgr.User{Email: gu.Email, LogoutUrl: lu}, err
}

// NewServer returns a new implementation for a UserServiceServer
func NewServer() finmgr.UserServiceServer {
	return new(server)
}

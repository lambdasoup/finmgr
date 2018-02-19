package user

import (
	"github.com/lambdasoup/finmgr"
	"github.com/lambdasoup/finmgr/aegrpc"
	"github.com/lambdasoup/finmgr/vapid"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	aeuser "google.golang.org/appengine/user"
)

type server struct{}

type user struct {
}

type subscription struct {
	Endpoint string
	P256dh   []byte
	Auth     []byte
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

func (sv *server) PutSubscription(ctx context.Context, in *finmgr.Subscription) (*finmgr.Empty, error) {
	actx := aegrpc.NewAppengineContext(ctx)
	uk := aegrpc.GetUserKey(actx)

	sk := datastore.NewIncompleteKey(actx, "Subscription", uk)
	s := subscription{Endpoint: in.GetEndpoint(), P256dh: in.GetP256Dh(), Auth: in.GetAuth()}
	_, err := datastore.Put(actx, sk, &s)

	if err == nil {
		vapid.SendTestMessageTo(actx, s.Endpoint, s.Auth, s.P256dh)
	}

	return &finmgr.Empty{}, err
}

// NewServer returns a new implementation for a UserServiceServer
func NewServer() finmgr.UserServiceServer {
	return new(server)
}

package push

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"net/http"

	"golang.org/x/net/context"

	webpush "github.com/SherClockHolmes/webpush-go"
	"github.com/lambdasoup/finmgr"
	"github.com/lambdasoup/finmgr/aegrpc"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

var curve = elliptic.P256()

type server struct{}

type subscription struct {
	Endpoint string
	P256dh   []byte
	Auth     []byte
}

type keyPair struct {
	X []byte
	Y []byte
	D []byte
}

// GetPublicKey handles a key request
func GetPublicKey(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)

	kk := datastore.NewKey(ctx, "KeyPair", "vapid-keypair", 0, nil)
	k := keyPair{}
	err := datastore.Get(ctx, kk, &k)

	// no key yet? let's build it now
	if err == datastore.ErrNoSuchEntity {
		key, err2 := ecdsa.GenerateKey(curve, rand.Reader)
		k.X = key.PublicKey.X.Bytes()
		k.Y = key.PublicKey.Y.Bytes()
		k.D = key.D.Bytes()
		datastore.Put(ctx, kk, &k)
		if err2 != nil {
			msg := fmt.Sprintf("could not generade key (%v)", err2)
			w.Write([]byte(msg))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		msg := fmt.Sprintf("could not load key (%v)", err)
		w.Write([]byte(msg))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bs := k.publicKey()
	w.Write(bs)
}

func (k keyPair) publicKey() []byte {
	// uncompressed pubkey format
	bs := []byte{0x04}
	bs = append(bs, k.X...)
	bs = append(bs, k.Y...)
	return bs
}

func sendTestMessageTo(ctx context.Context, ep string, auth []byte, p256dh []byte) {

	kk := datastore.NewKey(ctx, "KeyPair", "vapid-keypair", 0, nil)
	k := keyPair{}
	_ = datastore.Get(ctx, kk, &k)
	prvk := ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: curve,
			X:     new(big.Int),
			Y:     new(big.Int),
		},
		D: new(big.Int),
	}
	prvk.D.SetBytes(k.D)
	prvk.PublicKey.X.SetBytes(k.X)
	prvk.PublicKey.Y.SetBytes(k.Y)

	s := webpush.Subscription{}
	s.Endpoint = ep
	s.Keys = webpush.Keys{}
	s.Keys.Auth = base64.RawURLEncoding.EncodeToString(auth)
	s.Keys.P256dh = base64.RawURLEncoding.EncodeToString(p256dh)

	// Send Notification
	_, err := webpush.SendNotification([]byte("TestXYZ"), &s, &webpush.Options{
		Subscriber:      "<mh@lambdasoup.com>",
		VAPIDPrivateKey: base64.RawURLEncoding.EncodeToString(prvk.D.Bytes()),
		HTTPClient:      urlfetch.Client(ctx),
	})
	if err != nil {
		log.Errorf(ctx, "could not send notification: %v", err)
	}
}

func (sv *server) PutSubscription(ctx context.Context, in *finmgr.Subscription) (*finmgr.Empty, error) {
	actx := aegrpc.NewAppengineContext(ctx)
	uk := aegrpc.GetUserKey(ctx)

	sk := datastore.NewIncompleteKey(actx, "Subscription", uk)
	s := subscription{Endpoint: in.GetEndpoint(), P256dh: in.GetP256Dh(), Auth: in.GetAuth()}
	_, err := datastore.Put(actx, sk, &s)
	return &finmgr.Empty{}, err
}

// Notify sends the given message to the given user
func Notify(ctx context.Context, uk *datastore.Key, scope string) error {
	// get server keys
	kk := datastore.NewKey(ctx, "KeyPair", "vapid-keypair", 0, nil)
	k := keyPair{}
	_ = datastore.Get(ctx, kk, &k)
	prvk := ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: curve,
			X:     new(big.Int),
			Y:     new(big.Int),
		},
		D: new(big.Int),
	}
	prvk.D.SetBytes(k.D)
	prvk.PublicKey.X.SetBytes(k.X)
	prvk.PublicKey.Y.SetBytes(k.Y)

	// get user keys
	ss := []subscription{}
	_, err := datastore.NewQuery("Subscription").Ancestor(uk).GetAll(ctx, &ss)
	if err != nil {
		return err
	}

	// send pushes to each sub
	for _, s := range ss {
		ws := webpush.Subscription{}
		ws.Endpoint = s.Endpoint
		ws.Keys = webpush.Keys{}
		ws.Keys.Auth = base64.RawURLEncoding.EncodeToString(s.Auth)
		ws.Keys.P256dh = base64.RawURLEncoding.EncodeToString(s.P256dh)

		// Send Notification
		_, err = webpush.SendNotification([]byte(scope), &ws, &webpush.Options{
			Subscriber:      "<mh@lambdasoup.com>",
			VAPIDPrivateKey: base64.RawURLEncoding.EncodeToString(prvk.D.Bytes()),
			HTTPClient:      urlfetch.Client(ctx),
		})
		if err != nil {
			log.Warningf(ctx, "could not send notification: %v", err)
		}
	}

	return nil
}

// NewServer returns a new implementation for a UserServiceServer
func NewServer() finmgr.PushServiceServer {
	return new(server)
}

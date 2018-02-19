package vapid

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"net/http"

	webpush "github.com/SherClockHolmes/webpush-go"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

var curve = elliptic.P256()

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

// SendTestMessageTo sends a test message to the given subscriber
func SendTestMessageTo(ctx context.Context, ep string, auth []byte, p256dh []byte) {

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

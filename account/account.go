package account

import (
	"net/http"
	"net/url"

	"github.com/lambdasoup/finmgr"
	"github.com/lambdasoup/finmgr/aegrpc"
	"github.com/mitch000001/go-hbci/client"
	"github.com/mitch000001/go-hbci/domain"
	https "github.com/mitch000001/go-hbci/transport/https"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/taskqueue"
	"google.golang.org/appengine/urlfetch"
)

type server struct{}

type bank struct {
	ID      string
	BLZ     string
	PIN     string
	Loading bool
}

type account struct {
	Name string
}

func (s *server) AddBank(ctx context.Context, in *finmgr.Bank) (*finmgr.Empty, error) {
	actx := aegrpc.NewAppengineContext(ctx)
	uk := aegrpc.GetUserKey(ctx)

	//  schedule hbci work
	err := datastore.RunInTransaction(actx, func(tctx context.Context) error {
		b := bank{
			ID:      in.GetId(),
			BLZ:     in.GetBlz(),
			PIN:     in.GetPin(),
			Loading: true,
		}

		bk := makeBankKey(tctx, b.ID, uk)
		_, terr := datastore.Put(tctx, bk, &b)
		if terr != nil {
			return terr
		}

		vs := url.Values{}
		vs.Set("id", b.ID)
		vs.Set("uid", uk.StringID())
		t := taskqueue.NewPOSTTask("/worker", vs)
		_, terr = taskqueue.Add(tctx, t, "")
		if terr != nil {
			return terr
		}

		// TODO update client

		return nil
	}, nil)

	return &finmgr.Empty{}, err
}

func UpdateAccounts(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)

	// get bank
	id := req.FormValue("id")
	uid := req.FormValue("uid")
	uk := aegrpc.MakeUserKey(ctx, uid)
	bk := makeBankKey(ctx, id, uk)
	b := bank{}
	err := datastore.Get(ctx, bk, &b)
	if err != nil {
		log.Errorf(ctx, "could not get bank record: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// make HBCI query
	config := client.Config{
		AccountID:   b.ID,
		BankID:      b.BLZ,
		PIN:         b.PIN,
		HBCIVersion: domain.FINTSVersion300,
		Transport:   https.NewNonDefault(urlfetch.Client(ctx)),
	}

	c, err := client.New(config)
	if err != nil {
		log.Errorf(ctx, "could not create hbci client: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	has, err := c.Accounts()
	if err != nil {
		log.Errorf(ctx, "could not get hbci accounts: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	log.Infof(ctx, "loaded %v accounts", len(has))

	// TODO update db
	b.Loading = false
	_, err = datastore.Put(ctx, bk, &b)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// TODO update client
	// vapid.UpdateClient(ctx, uk, "AccountService")
}

func (s *server) GetAccounts(ctx context.Context, in *finmgr.Empty) (*finmgr.Accounts, error) {
	actx := aegrpc.NewAppengineContext(ctx)
	uk := aegrpc.GetUserKey(ctx)

	bs := []bank{}
	bks, err := datastore.NewQuery("Bank").Ancestor(uk).KeysOnly().GetAll(actx, &bs)
	if err != nil {
		return nil, err
	}

	// TODO limited to one bank for now
	if len(bks) == 0 {
		return &finmgr.Accounts{Loading: false, List: []*finmgr.Account{}}, nil
	}

	as := []account{}
	_, err = datastore.NewQuery("Account").Ancestor(bks[0]).GetAll(actx, &as)
	if err != nil {
		return nil, err
	}

	pbas := make([]*finmgr.Account, len(as))
	for i := range as {
		pbas[i].Name = as[i].Name
	}

	return &finmgr.Accounts{Loading: bs[0].Loading, List: pbas}, nil
}

func makeBankKey(ctx context.Context, id string, uk *datastore.Key) *datastore.Key {
	return datastore.NewKey(ctx, "Bank", id, 0, uk)
}

// NewServer returns a new gRPC accound handler
func NewServer() finmgr.AccountServiceServer {
	return new(server)
}

package account

import (
	"net/http"
	"net/url"
	"time"

	"github.com/lambdasoup/finmgr"
	"github.com/lambdasoup/finmgr/aegrpc"
	"github.com/lambdasoup/finmgr/push"
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
	ID  string
	BLZ string
	PIN string
	// TODO get rid of this?
	Loading bool
}

type account struct {
	Name string
}

// AddBank adds a new bank with the given arguments
func (s *server) AddBank(ctx context.Context, in *finmgr.AddBankRequest) (*finmgr.Empty, error) {
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

		// update client
		push.Notify(actx, uk, "Account")

		return nil
	}, nil)

	return &finmgr.Empty{}, err
}

// RefreshBank refreshed the given bank's accounts
func (s *server) RefreshBank(ctx context.Context, in *finmgr.RefreshRequest) (*finmgr.Empty, error) {
	actx := aegrpc.NewAppengineContext(ctx)
	uk := aegrpc.GetUserKey(ctx)

	//  schedule hbci work
	err := datastore.RunInTransaction(actx, func(tctx context.Context) error {
		// get bank
		b := bank{}
		bk := makeBankKey(tctx, in.GetId(), uk)
		terr := datastore.Get(tctx, bk, &b)
		if terr != nil {
			return terr
		}

		// set loading
		b.Loading = true
		_, terr = datastore.Put(tctx, bk, &b)
		if terr != nil {
			return terr
		}

		// remove current accounts
		aks, terr := datastore.NewQuery("Account").
			Ancestor(bk).
			KeysOnly().
			GetAll(tctx, nil)
		if terr != nil {
			return terr
		}
		terr = datastore.DeleteMulti(tctx, aks)
		if terr != nil {
			return terr
		}

		// enqueue task
		vs := url.Values{}
		vs.Set("id", b.ID)
		vs.Set("uid", uk.StringID())
		t := taskqueue.NewPOSTTask("/worker/update-accounts", vs)
		_, terr = taskqueue.Add(tctx, t, "")
		if terr != nil {
			return terr
		}

		return nil
	}, nil)

	// update client
	push.Notify(actx, uk, "Account")

	return &finmgr.Empty{}, err
}

// UpdateAccounts updates accounts for the given bank
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

	// update bank entity
	err = datastore.RunInTransaction(ctx, func(tctx context.Context) error {
		// insert accounts into db
		for _, ha := range has {
			a := account{
				Name: ha.ProductID,
			}
			ak := datastore.NewKey(ctx, "Account", ha.AccountConnection.AccountID, 0, bk)
			_, terr := datastore.Put(ctx, ak, &a)
			if terr != nil {
				return terr
			}
		}

		// update bank loading state
		b.Loading = false
		_, terr := datastore.Put(ctx, bk, &b)
		return terr
	}, nil)
	if err != nil {
		log.Errorf(ctx, "could update bank record: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// update client
	push.Notify(ctx, uk, "Account")
}

// UpdateTransactions updates transactions for all banks
func UpdateTransactions(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)

	// get banks
	bs := []bank{}
	bks, err := datastore.NewQuery("Bank").GetAll(ctx, &bs)
	if err != nil {
		log.Errorf(ctx, "could not get bank record: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if len(bs) != 1 {
		log.Errorf(ctx, "wanted exactly 1 bank, but got %v", len(bs))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	b := bs[0]
	bk := bks[0]
	uk := bk.Parent()

	// make HBCI query
	config := client.Config{
		AccountID:   b.ID,
		BankID:      b.BLZ,
		PIN:         b.PIN,
		HBCIVersion: domain.FINTSVersion300,
		Transport:   https.NewNonDefault(urlfetch.Client(ctx)),
	}
	log.Debugf(ctx, "built hbci config: %v", config)

	c, err := client.New(config)
	if err != nil {
		log.Errorf(ctx, "could not create hbci client: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// get accounts TODO from db/hbci/both?
	// as := []account{}
	// aks, err := datastore.NewQuery("Account").GetAll(ctx, &as)
	// if err != nil {
	// 	log.Errorf(ctx, "could not get accounts for bank: %v", err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	w.Write([]byte(err.Error()))
	// 	return
	// }
	has, err := c.Accounts()
	if err != nil {
		log.Errorf(ctx, "could not get hbci accounts: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	d := time.Date(2017, 1, 1, 0, 0, 0, 0, time.Local)
	sd := domain.NewShortDate(d)
	tf := domain.TimeframeFromDate(sd)
	txs, err := c.AccountTransactions(has[0].AccountConnection, tf, false, "")
	if err != nil {
		log.Errorf(ctx, "could not get hbci transactions: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	log.Infof(ctx, "got txs: %v", txs)

	// update client
	push.Notify(ctx, uk, "Account")
}

func (s *server) GetBanks(ctx context.Context, in *finmgr.Empty) (*finmgr.BanksResponse, error) {
	actx := aegrpc.NewAppengineContext(ctx)
	uk := aegrpc.GetUserKey(ctx)

	bs := []bank{}
	bks, err := datastore.NewQuery("Bank").Ancestor(uk).GetAll(actx, &bs)
	if err != nil {
		return nil, err
	}

	resp := finmgr.BanksResponse{}
	for _, b := range bs {
		pbb := finmgr.Bank{}
		pbb.Blz = b.BLZ
		pbb.Id = b.ID
		pbb.Updating = b.Loading

		as := []account{}
		aks, err := datastore.NewQuery("Account").Ancestor(bks[0]).GetAll(actx, &as)
		if err != nil {
			return nil, err
		}
		for i, a := range as {
			pbb.Accounts = append(pbb.Accounts, &finmgr.Account{Name: a.Name, Id: aks[i].StringID()})
		}

		resp.Banks = append(resp.Banks, &pbb)
	}

	return &resp, nil
}

func makeBankKey(ctx context.Context, id string, uk *datastore.Key) *datastore.Key {
	return datastore.NewKey(ctx, "Bank", id, 0, uk)
}

// NewServer returns a new gRPC accound handler
func NewServer() finmgr.AccountServiceServer {
	return new(server)
}

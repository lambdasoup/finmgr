package account

import (
	"fmt"

	"github.com/lambdasoup/finmgr"
	"github.com/lambdasoup/finmgr/aegrpc"
	"github.com/mitch000001/go-hbci/client"
	"github.com/mitch000001/go-hbci/domain"
	https "github.com/mitch000001/go-hbci/transport/https"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

type server struct{}

func (s *server) GetAccounts(ctx context.Context, in *finmgr.AccountInfo) (*finmgr.Accounts, error) {
	actx := aegrpc.NewAppengineContext(ctx)
	log.Debugf(actx, "in %v", in)

	config := client.Config{
		AccountID:   in.GetId(),
		BankID:      in.GetBlz(),
		PIN:         in.GetPin(),
		HBCIVersion: domain.FINTSVersion300,
		Transport:   https.NewNonDefault(urlfetch.Client(actx)),
	}

	c, err := client.New(config)
	if err != nil {
		panic(err)
	}

	accounts, err := c.Accounts()

	return &finmgr.Accounts{Info: fmt.Sprintf("TODO accounts info %v", accounts)}, err
}

func NewServer() finmgr.AccountServiceServer {
	return new(server)
}

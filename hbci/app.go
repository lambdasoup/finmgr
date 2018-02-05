package ui

import (
	"fmt"
	"net/http"

	"github.com/mitch000001/go-hbci/client"
	"github.com/mitch000001/go-hbci/domain"
	https "github.com/mitch000001/go-hbci/transport/https"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

func init() {
	http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	config := client.Config{
		URL:         "AAA",
		AccountID:   "XXX",
		BankID:      "YYY",
		PIN:         "ZZZ",
		HBCIVersion: domain.FINTSVersion300,
		Transport:   https.NewNonDefault(urlfetch.Client(ctx)),
	}

	c, err := client.New(config)
	if err != nil {
		panic(err)
	}

	accounts, err := c.Accounts()

	if err != nil {
		log.Debugf(ctx, "cannot get accounts: %v", err)
	}

	fmt.Fprintf(w, "%v/n", accounts)
	fmt.Fprint(w, "Hello, UI!")
}

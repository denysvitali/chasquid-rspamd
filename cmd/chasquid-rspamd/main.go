package main

import (
	"fmt"
	"github.com/alexflint/go-arg"
	rspamd "github.com/denysvitali/chasquid-rspamd/pkg"
	"net/url"
	"os"
)
import "github.com/sirupsen/logrus"

var args struct {
	RspamdURL string `arg:"-u,--url" help:"rspamd URL"`
	Auth      string `arg:"-a,--auth" help:"Auth string"`
}

func main() {
	logger := logrus.New()
	logger.Out = os.Stderr
	arg.MustParse(&args)
	u, err := url.Parse(args.RspamdURL)
	if err != nil {
		logger.Fatalf("invalid rspamd URL: %v", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		logger.Fatalf("invalid scheme, only http and https are supported")
	}

	c, err := rspamd.New(args.RspamdURL)
	if err != nil {
		logger.Fatalf("unable to create rspamd client: %v", err)
	}
	if args.Auth != "" {
		c.SetAuth(args.Auth)
	}

	res, err := c.Scan(&rspamd.ScanRequest{Body: os.Stdin})
	if err != nil {
		logger.Fatalf("unable to perform scan: %v", err)
	}

	fmt.Printf("X-Spam-Action: %s\n", res.Action)
	fmt.Printf("X-Spam-Score: %.2f\n", res.Score)
}

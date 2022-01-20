package main

import (
	"github.com/alexflint/go-arg"
	"net/url"
)
import "github.com/sirupsen/logrus"

var args struct {
	RspamdURL string `arg:"-u,--url" help:"rspamd URL"`
}

func main() {
	logger := logrus.New()
	arg.MustParse(&args)
	u, err := url.Parse(args.RspamdURL)
	if err != nil {
		logger.Fatalf("invalid rspamd URL: %v", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		logger.Fatalf("invalid scheme, only http and https are supported")
	}
}

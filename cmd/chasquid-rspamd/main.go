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
	LogToFile string `arg:"-f,--log-file" help:"Log to file"`
	Debug     *bool  `arg:"-D,--debug"`
}

func main() {
	logger := logrus.New()
	arg.MustParse(&args)

	var outputLogFile *os.File = nil
	var err error
	defer func() {
		if outputLogFile != nil {
			outputLogFile.Close()
		}
	}()

	if args.LogToFile != "" {
		outputLogFile, err = os.OpenFile(args.LogToFile,
			os.O_CREATE | os.O_APPEND | os.O_RDWR,
			0644,
		)
		if err != nil {
			logger.Fatalf("unable to set output log file: %v", err)
		}
		logger.Out = outputLogFile
	} else {
		logger.Out = os.Stderr
	}

	if args.Debug != nil && *args.Debug {
		logger.SetLevel(logrus.DebugLevel)
	}

	u, err := url.Parse(args.RspamdURL)
	if err != nil {
		logger.Fatalf("invalid rspamd URL: %v", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		logger.Fatalf("invalid scheme %s, only http and https are supported", u.Scheme)
	}

	c, err := rspamd.New(args.RspamdURL)
	if err != nil {
		logger.Fatalf("unable to create rspamd client: %v", err)
	}

	req := rspamd.ScanRequest{
		Body:     os.Stdin,
		SourceIP: envVarOrNil("REMOTE_ADDR"),
		From:     envVarOrNil("MAIL_FROM"),
		Hostname: envVarOrNil("EHLO_DOMAIN"),
		Rcpt:     envVarOrNil("RCPT_TO"),
	}
	if args.Auth != "" {
		req.Password = &args.Auth
	}

	res, err := c.Scan(&req)
	if err != nil {
		logger.Fatalf("unable to perform scan: %v", err)
	}

	logger.Debugf("response=%v", res)

	fmt.Printf("X-Spam-Action: %s\n", res.Action)
	fmt.Printf("X-Spam-Score: %.2f\n", res.Score)
}

func envVarOrNil(key string) *string {
	v := os.Getenv(key)
	if v == "" {
		return nil
	}
	return &v
}
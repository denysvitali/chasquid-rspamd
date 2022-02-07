package rspamd

import (
	"fmt"
	"net/http"
	"net/url"
)

type HttpClient interface {
	Do(r *http.Request) (*http.Response, error)
}

type Client struct {
	HttpClient HttpClient
	RspamdUrl  *url.URL
}

func New(rspamdUrl string) (*Client, error) {
	u, err := url.Parse(rspamdUrl)
	if err != nil {
		return nil, fmt.Errorf("unable to parse rspamd URL: %v", err)
	}

	c := Client{
		HttpClient: http.DefaultClient,
		RspamdUrl:  u,
	}

	return &c, nil
}
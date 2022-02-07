package rspamd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type PassType = string

type ScanRequest struct {
	Body          io.ReadCloser
	DeliverTo     *string
	SourceIP      *string
	Helo          *string
	Hostname      *string
	Flags         *string
	From          *string
	QueueId       *string
	Raw           *bool
	Rcpt          *string
	Pass          *string // If this header has all value, all filters would be checked for this message.
	Subject       *string
	User          *string
	MessageLength *int
	SettingsId    *string
	Settings      *string
	UserAgent     *string
	MTATag        *string
	MTAName       *string
	TLSCipher     *string
	TLSVersion    *string
	TLSCertIssuer *string
	UrlFormat     *string
	Filename      *string
}

type SymbolResult struct {
	Name  string
	Score float32
}

type ScanResult struct {
	IsSkipped     bool                    `json:"is_skipped"`
	Score         float64                 `json:"score"`
	RequiredScore float64                 `json:"required_score"`
	Action        string                  `json:"action"`
	Symbols       map[string]SymbolResult `json:"symbols"`
	Urls          []string                `json:"urls"`
	Emails        []string                `json:"emails"`
	MessageId     string                  `json:"message-id"`
}

func (c *Client) Scan(scanRequest *ScanRequest) (*ScanResult, error) {
	if scanRequest == nil {
		return nil, fmt.Errorf("invalid scanRequest")
	}

	checkV2, err := c.RspamdUrl.Parse("/checkv2")
	if err != nil {
		return nil, fmt.Errorf("unable to parse URL: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, checkV2.String(), scanRequest.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %v", err)
	}

	addEntry(req, "Deliver-To", scanRequest.DeliverTo)
	addEntry(req, "IP", scanRequest.SourceIP)
	addEntry(req, "Helo", scanRequest.Helo)
	addEntry(req, "Hostname", scanRequest.Hostname)
	addEntry(req, "Flags", scanRequest.Flags)
	addEntry(req, "From", scanRequest.From)
	addEntry(req, "Queue-Id", scanRequest.QueueId)
	addBoolEntry(req, "Raw", scanRequest.Raw)
	addEntry(req, "Rcpt", scanRequest.Rcpt)
	addEntry(req, "Pass", scanRequest.Pass)
	addEntry(req, "Subject", scanRequest.Subject)
	addEntry(req, "User", scanRequest.User)
	addIntEntry(req, "Message-Length", scanRequest.MessageLength)
	addEntry(req, "Settings-Id", scanRequest.SettingsId)
	addEntry(req, "Settings", scanRequest.Settings)
	addEntry(req, "User-Agent", scanRequest.UserAgent)
	addEntry(req, "MTA-Tag", scanRequest.MTATag)
	addEntry(req, "MTA-Name", scanRequest.MTAName)
	addEntry(req, "TLS-Cipher", scanRequest.TLSCipher)
	addEntry(req, "TLS-Version", scanRequest.TLSVersion)
	addEntry(req, "TLS-Cert-Issuer", scanRequest.TLSCertIssuer)
	addEntry(req, "URL-Format", scanRequest.UrlFormat)
	addEntry(req, "Filename", scanRequest.UrlFormat)

	res, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to execute request: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %s, 200 OK expected", res.Status)
	}

	var scanResult ScanResult
	d := json.NewDecoder(res.Body)
	err = d.Decode(&scanResult)
	if err != nil {
		return nil, fmt.Errorf("unable to decode JSON: %v", err)
	}
	return &scanResult, nil
}

func addEntry(req *http.Request, key string, v *string) {
	if v != nil {
		req.Header.Add(key, *v)
	}
}

func addBoolEntry(req *http.Request, key string, v *bool) {
	if v != nil {
		req.Header.Add(key, fmt.Sprintf("%v", *v))
	}
}

func addIntEntry(req *http.Request, key string, v *int) {
	if v != nil {
		req.Header.Add(key, fmt.Sprintf("%d", *v))
	}
}
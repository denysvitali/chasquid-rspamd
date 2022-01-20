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
	Hostname      *string
	Flags         *string
	From          *string
	QueueId       *string
	Raw           *bool
	Pass          *string // If this header has all value, all filters would be checked for this message.
	Subject       *string
	User          *string
	MessageLength *int
	SettingsId    *string
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

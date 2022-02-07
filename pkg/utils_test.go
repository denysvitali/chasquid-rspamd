package rspamd_test

import (
	"fmt"
	rspamd "github.com/denysvitali/chasquid-rspamd/pkg"
	"net/http"
)

type MockedHTTPClient struct {
	responseArray []http.Response
}

func (m *MockedHTTPClient) Do(r *http.Request) (*http.Response, error) {
	if len(m.responseArray) == 0 {
		return nil, fmt.Errorf("no more responses")
	}

	resp := m.responseArray[0]
	m.responseArray = m.responseArray[1:]
	return &resp, nil
}

var _ rspamd.HttpClient = (*MockedHTTPClient)(nil)

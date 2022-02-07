package rspamd_test

import (
	"fmt"
	rspamd "github.com/denysvitali/chasquid-rspamd/pkg"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

func TestClient(t *testing.T) {
	c, err := rspamd.New("http://127.0.0.1:11333")
	assert.Nil(t, err)
	assert.NotNil(t, c)

	f, err := os.Open("../sample/1.msg")
	if err != nil {
		t.Fatalf("unable to open file: %v", err)
	}

	fromEmail := "admin@tfujiyama.com"
	fromHostname := "f10.my.com"
	fromIp := "185.30.176.240"
	subject := "=?UTF-8?B?MTEgMCAg?="

	res, err := c.Scan(&rspamd.ScanRequest{
		Body:     f,
		SourceIP: &fromIp,
		Hostname: &fromHostname,
		From:     &fromEmail,
		Subject:  &subject,
	})
	assert.Nil(t, err)
	assert.NotNil(t, res)
	fmt.Printf("%+v\n", res)
	fmt.Printf("Score: %v\n", res.Score)
	for k, v := range res.Symbols {
		fmt.Printf("%s: %.2f\n", k, v.Score)
	}
	assert.Equal(t, "", res.Action)
}
func TestMockedClient(t *testing.T) {
	c, err := rspamd.New("http://127.0.0.1:11333")
	mockedHttpClient := MockedHTTPClient{}
	c.HttpClient = &mockedHttpClient

	f, err := os.Open("../sample/1.json")
	if err != nil {
		t.Fatalf("unable to open file: %v", err)
	}

	mockedHttpClient.responseArray = []http.Response{
		{Body: f, StatusCode: http.StatusOK},
	}

	assert.Nil(t, err)
	assert.NotNil(t, c)

	f, err = os.Open("../sample/1.msg")
	if err != nil {
		t.Fatalf("unable to open file: %v", err)
	}

	fromEmail := "admin@tfujiyama.com"
	fromHostname := "f10.my.com"
	fromIp := "185.30.176.240"
	subject := "=?UTF-8?B?MTEgMCAg?="

	res, err := c.Scan(&rspamd.ScanRequest{
		Body:     f,
		SourceIP: &fromIp,
		Hostname: &fromHostname,
		From:     &fromEmail,
		Subject:  &subject,
	})
	assert.Nil(t, err)
	assert.NotNil(t, res)

	assert.Equal(t, "reject", res.Action)
	assert.Equal(t, false, res.IsSkipped)
	assert.Equal(t, "1642699320.209195217@f10.my.com", res.MessageId)
	assert.Equal(t, 7.0, res.RequiredScore)
	assert.Equal(t, float32(1.0), res.Symbols["MIME_BASE64_TEXT_BOGUS"].Score)
}
func TestFailedRequest(t *testing.T) {
	c, err := rspamd.New("http://127.0.0.1:11333")
	mHttpClient := MockedHTTPClient{}
	c.HttpClient = &mHttpClient

	mHttpClient.responseArray = []http.Response{
		{StatusCode: http.StatusInternalServerError, Body: nil},
	}

	f, err := os.Open("../sample/1.msg")
	if err != nil {
		t.Fatalf("unable to open file: %v", err)
	}

	res, err := c.Scan(&rspamd.ScanRequest{
		Body: f,
	})
	assert.Nil(t, res)
	assert.NotNil(t, err)
}

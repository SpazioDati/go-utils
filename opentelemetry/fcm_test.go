package opentelemetry

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"net/http"
	"testing"
)

// Mock token generation
type mockTokenSource struct {
}

func (ts mockTokenSource) Token() (*oauth2.Token, error) {
	return &oauth2.Token{
		AccessToken: "test",
	}, nil
}
func mockGenerateTokenProvider([]byte) (*tokenProvider, error) {
	return &tokenProvider{tokenSource: mockTokenSource{}}, nil
}

// Mock http requests with a custom RoundTripper
type mockRoundTripper struct {
}

func (rt mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		Header: req.Header,
	}, nil
}

func TestFCMRequestHeaders(t *testing.T) {
	var c http.Client
	c.Transport = CustomFCMTransport(mockRoundTripper{}, mockGenerateTokenProvider, nil)
	res, err := c.Get("test")

	assert.NoError(t, err)
	assert.Equal(t, "Bearer test", res.Header.Get("Authorization"))
}

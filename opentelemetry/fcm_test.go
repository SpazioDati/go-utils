package opentelemetry

import (
	"errors"
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
func mockGenerateTokenProviderError([]byte) (*tokenProvider, error) {
	return nil, errors.New("token test error")
}

// Mock http requests with a custom RoundTripper
type mockRoundTripper struct {
	sendError bool
}
func (rt mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if rt.sendError {
		return nil, errors.New("test error")
	}
	return &http.Response{
		Header: req.Header,
	}, nil
}

func TestFCMRequestHeaders(t *testing.T) {
	var c http.Client
	c.Transport = CustomFCMTransport(mockRoundTripper{sendError: false}, mockGenerateTokenProvider, nil)
	res, err := c.Get("test")

	assert.NoError(t, err)
	assert.Equal(t, "Bearer test", res.Header.Get("Authorization"))
}

func TestFCMRequestWithError(t *testing.T) {
	var c http.Client
	c.Transport = CustomFCMTransport(mockRoundTripper{sendError: true}, mockGenerateTokenProvider, nil)
	_, err := c.Get("test")

	assert.Error(t, err)
	assert.Equal(t, "Get \"test\": test error", err.Error())
}

func TestFCMRequestWithTokenError(t *testing.T) {
	var c http.Client
	c.Transport = CustomFCMTransport(mockRoundTripper{sendError: false}, mockGenerateTokenProviderError, nil)
	_, err := c.Get("test")

	assert.Error(t, err)
	assert.Equal(t, "Get \"test\": token test error", err.Error())
}

func TestNoRoundtripperProvided(t *testing.T) {
	customTransport := CustomFCMTransport(nil, mockGenerateTokenProvider, nil)

	assert.NotNil(t, customTransport.T)
	assert.Equal(t, customTransport.T, http.DefaultTransport)
}

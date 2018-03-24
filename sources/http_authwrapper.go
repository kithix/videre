package sources

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bobziuchkovski/digest"
)

// This could possible use net.URL
type HttpRequestDetails struct {
	URL      string
	AuthType string
	User     string
	Pass     string
}

var defaultClient = &http.Client{
	Timeout: 5 * time.Second,
}

func buildNoAuthRequest(url string) (*http.Request, *http.Client, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	return request, defaultClient, nil
}

func buildBasicAuthRequest(url, user, pass string) (*http.Request, *http.Client, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	request.SetBasicAuth(user, pass)

	return request, defaultClient, nil
}

func buildDigestAuthRequest(url, user, pass string) (*http.Request, *http.Client, error) {
	transport := digest.NewTransport(user, pass)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	client, err := transport.Client()
	if err != nil {
		return nil, nil, err
	}

	return request, client, nil
}

func buildRequest(d HttpRequestDetails) (*http.Request, *http.Client, error) {
	authType := strings.ToLower(d.AuthType)
	// Factory implementation for auth
	if authType == "" {
		return buildNoAuthRequest(d.URL)
	}
	if authType == "basic" {
		return buildBasicAuthRequest(d.URL, d.User, d.Pass)
	}
	if authType == "digest" {
		return buildDigestAuthRequest(d.URL, d.User, d.Pass)
	}

	return nil, nil, fmt.Errorf("No request maker for AuthType: " + authType + " available")
}

package redmine

import (
	"crypto/tls"
	"github.com/pkg/errors"
	"net/http"
)

// ClientBuilder takes configurations and safely builds a Redmine client instance.
//
// Minimal example:
//  client, err := NewClientBuilder().
//    Endpoint("https://localhost:3000").
//    AuthBasicAuth("adminUser", "passwort1"). // or AuthAPIToken(apiToken)
//    Build()
//
// Example for Basic Authentication, skipping self-signed SSL certificates:
//  client, err := NewClientBuilder().
//    Endpoint("https://localhost:3000").
//    AuthBasicAuth("adminUser", "passwort1").
//    SkipSSLVerify(true).
//    Build()
type ClientBuilder struct {
	endpoint      string
	auth          *APIAuth
	skipSSLVerify bool
	resultLimit   int
	resultOffset  int
}

// NewClientBuilder creates a new ClientBuilder which in turn .
func NewClientBuilder() *ClientBuilder {
	return &ClientBuilder{
		resultLimit:  NoSetting,
		resultOffset: NoSetting,
	}
}

// Endpoint configures the URL under which the Redmine is available.
func (cb *ClientBuilder) Endpoint(endpoint string) *ClientBuilder {
	cb.endpoint = endpoint
	return cb
}

// AuthBasicAuth configures an APIAuth instance to use basic authentication credentials. At least one way of authentication must be configured.
func (cb *ClientBuilder) AuthBasicAuth(user, password string) *ClientBuilder {
	cb.auth = &APIAuth{AuthType: AuthTypeBasicAuth, User: user, Password: password}
	return cb
}

// AuthBasicAuth configures an APIAuth instance to use API query tokens. At least one way of authentication must be configured.
func (cb *ClientBuilder) AuthAPIToken(apiKey string) *ClientBuilder {
	cb.auth = &APIAuth{AuthType: AuthTypeTokenQueryParam, Token: apiKey}
	return cb
}

// SkipSSLVerify (if set to true) configures the HTTP client to skip SSL certificate verification.
func (cb *ClientBuilder) SkipSSLVerify(skip bool) *ClientBuilder {
	cb.skipSSLVerify = skip
	return cb
}

// ResultLimit sets the limit of elements in a paged result.
func (cb *ClientBuilder) ResultLimit(limit int) *ClientBuilder {
	cb.resultLimit = limit
	return cb
}

// ResultOffset sets the offset for result paging.
func (cb *ClientBuilder) ResultOffset(offset int) *ClientBuilder {
	cb.resultOffset = offset
	return cb
}

// Build validates the given client configuration and returns a client unless there are configuration errors.
func (cb *ClientBuilder) Build() (*Client, error) {
	var err error
	if err = cb.assertEndpoint(); err != nil {
		return nil, errors.Wrapf(err, "could not create redmine client")
	}

	if err = cb.assertAuth(); err != nil {
		return nil, errors.Wrapf(err, "could not create redmine client")
	}

	httpClient := http.DefaultClient
	if cb.skipSSLVerify {
		httpClient = skipSslHttpClient()
	}

	return &Client{
		endpoint: cb.endpoint,
		auth:     *cb.auth,
		Limit:    cb.resultLimit,
		Offset:   cb.resultOffset,
		Client:   httpClient,
	}, nil
}

func (cb *ClientBuilder) assertAuth() error {
	if cb.auth == nil {
		return errors.New("redmine authentication must not be empty")
	}

	err := cb.auth.validate()
	if err != nil {
		return err
	}

	return nil
}

func (cb *ClientBuilder) assertEndpoint() error {
	if cb.endpoint == "" {
		return errors.New("redmine endpoint must not be empty")
	}

	return nil
}

func skipSslHttpClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return &http.Client{Transport: tr}
}

package redmine

import (
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Client struct {
	endpoint string
	auth     APIAuth
	Limit    int
	Offset   int
	*http.Client
}

const NoSetting = -1
const (
	AuthTypeBasicAuth = iota
	AuthTypeTokenQueryParam
	AuthTypeBasicAuthWithTokenPassword
	AuthTypeNoAuth
)

var DefaultLimit int = NoSetting
var DefaultOffset int = NoSetting

type AuthType int

type APIAuth struct {
	AuthType AuthType
	Token    string
	User     string
	Password string
}

func (auth APIAuth) validate() error {
	if auth.AuthType < AuthTypeBasicAuth || auth.AuthType > AuthTypeNoAuth {
		return fmt.Errorf("invalid auth configuration: AuthType %d found", auth.AuthType)
	}

	if auth.AuthType == AuthTypeBasicAuth || auth.AuthType == AuthTypeBasicAuthWithTokenPassword {
		if auth.User == "" {
			return fmt.Errorf("invalid auth configuration for type %d: user must not be empty", auth.AuthType)
		}
	}

	if auth.AuthType == AuthTypeTokenQueryParam || auth.AuthType == AuthTypeBasicAuthWithTokenPassword {
		if auth.Token == "" {
			return fmt.Errorf("invalid auth configuration for type %d: API token must not be empty", auth.AuthType)
		}
	}

	return nil
}

func NewClient(endpoint string, auth APIAuth) (*Client, error) {
	if err := auth.validate(); err != nil {
		return nil, errors.Wrapf(err, "could not create redmine client")
	}
	client := &Client{
		endpoint: endpoint,
		auth:     auth,
		Limit:    DefaultLimit,
		Offset:   DefaultOffset,
		Client:   http.DefaultClient,
	}

	return client, nil
}

func (c *Client) authenticatedGet(urlWithoutAuthInfo string) (req *http.Request, err error) {
	const method = "GET"
	errorMsg := fmt.Sprintf("could not create %s request for %s and auth type %d", method, urlWithoutAuthInfo, c.auth.AuthType)

	req, err = http.NewRequest(method, urlWithoutAuthInfo, nil)

	switch c.auth.AuthType {
	case AuthTypeBasicAuth:
		if err != nil {
			return nil, errors.Wrap(err, errorMsg)
		}
		req.SetBasicAuth(c.auth.User, c.auth.Password)
		return req, nil
	case AuthTypeTokenQueryParam:
		modifiedURL, err := safelyAddQueryParameter(urlWithoutAuthInfo, "key", c.auth.Token)
		if err != nil {
			return nil, errors.Wrap(err, errorMsg)
		}
		req, err = http.NewRequest(method, modifiedURL, nil)
		if err != nil {
			return nil, errors.Wrap(err, errorMsg)
		}
		return req, nil
	case AuthTypeBasicAuthWithTokenPassword:
		if err != nil {
			return nil, errors.Wrap(err, errorMsg)
		}
		req.SetBasicAuth(c.auth.User, c.auth.Token)
		return req, nil
	case AuthTypeNoAuth:
		if err != nil {
			return nil, errors.Wrap(err, errorMsg)
		}
		return req, nil
	}

	return nil, errors.New("unsupported auth type") // must never occur because it was validated earlier
}

func safelyAddQueryParameter(theURL string, key, value string) (string, error) {
	parsedURL, err := url.Parse(theURL)
	if err != nil {
		return "", errors.Wrapf(err, "could not add query parameter %s because parsing the URL %s failed", key, theURL)
	}
	parsedURL.Query().Add(key, value)
	return parsedURL.String(), nil
}

func (c *Client) apiKeyParameter() string {
	return "key=" + c.auth.Token
}

func (c *Client) concatParameters(requestParameters ...string) string {
	cleanedParams := []string{}
	for _, param := range requestParameters {
		if param != "" {
			cleanedParams = append(cleanedParams, param)
		}
	}

	return strings.Join(cleanedParams, "&")
}

// URLWithFilter return string url by concat endpoint, path and filter
// err != nil when endpoint can not parse
func (c *Client) URLWithFilter(path string, f Filter) (string, error) {
	var fullURL *url.URL
	fullURL, err := url.Parse(c.endpoint)
	if err != nil {
		return "", err
	}
	fullURL.Path += path
	if c.Limit > -1 {
		f.AddPair("limit", strconv.Itoa(c.Limit))
	}
	if c.Offset > -1 {
		f.AddPair("offset", strconv.Itoa(c.Offset))
	}
	fullURL.RawQuery = f.ToURLParams()
	return fullURL.String(), nil
}

func (c *Client) getPaginationClause() string {
	clause := ""
	if c.Limit > -1 {
		clause = clause + fmt.Sprintf("&limit=%d", c.Limit)
	}
	if c.Offset > -1 {
		clause = clause + fmt.Sprintf("&offset=%d", c.Offset)
	}
	return clause
}

type errorsResult struct {
	Errors []string `json:"errors"`
}

type IdName struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Id struct {
	Id int `json:"id"`
}

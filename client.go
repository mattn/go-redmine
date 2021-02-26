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
		return fmt.Errorf("invalid AuthType %d found", auth.AuthType)
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

func (c *Client) buildAuthenticatedURL(urlWithoutAuthInfo string) string {
	switch c.auth.AuthType {

	}
	return ""
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

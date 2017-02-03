package redmine

import (
	"fmt"
	"net/http"
)

type Client struct {
	endpoint string
	apikey   string
	*http.Client
	Limit  int
	Offset int
}

var DefaultLimit int = -1  // "-1" means "No setting"
var DefaultOffset int = -1 //"-1" means "No setting"

func NewClient(endpoint, apikey string) *Client {
	return &Client{endpoint, apikey, http.DefaultClient, DefaultLimit, DefaultOffset}
}
func (c *Client) getPaginationClause() string {
	clause := ""
	if c.Limit > -1 {
		clause = clause + fmt.Sprintf("&limit=%v", c.Limit)
	}
	if c.Offset > -1 {
		clause = clause + fmt.Sprintf("&offset=%v", c.Offset)
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

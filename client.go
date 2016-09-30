package redmine

import (
	"fmt"
	"net/http"
)

type client struct {
	endpoint string
	apikey   string
	*http.Client
	Limit  int
	Offset int
}

var DefaultLimit int = -1  // "-1" means "No setting"
var DefaultOffset int = -1 //"-1" means "No setting"

func NewClient(endpoint, apikey string) *client {
	return &client{endpoint, apikey, http.DefaultClient, DefaultLimit, DefaultOffset}
}
func (c *client) getPaginationClause() string {
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

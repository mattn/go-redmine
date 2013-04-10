package redmine

import "net/http"

type client struct {
	endpoint string
	apikey   string
	*http.Client
}

func NewClient(endpoint, apikey string) *client {
	return &client{endpoint, apikey, http.DefaultClient}
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

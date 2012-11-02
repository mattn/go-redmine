package redmine

type client struct {
	endpoint string
	apikey string
}

func NewClient(endpoint, apikey string) *client {
	return &client { endpoint, apikey }
}

type errorsResult struct {
	Errors []string `json:"errors"`
}

type IdName struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

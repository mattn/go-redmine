package redmine

import (
	"encoding/json"
	"errors"
	"strings"
)

type rolesResult struct {
	Roles []IdName `json:"roles"`
}

func (c *client) Roles() ([]IdName, error) {
	res, err := c.Get(c.endpoint + "/roles.json?key=" + c.apikey)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r rolesResult
	if res.StatusCode != 200 {
		var er errorsResult
		err = decoder.Decode(&er)
		if err == nil {
			err = errors.New(strings.Join(er.Errors, "\n"))
		}
	} else {
		err = decoder.Decode(&r)
	}
	if err != nil {
		return nil, err
	}
	return r.Roles, nil
}

package redmine

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type customFieldsResult struct {
	CustomFields []CustomField `json:"custom_fields"`
}

// CustomFields consulta los campos personalizados
func (c *Client) CustomFields() ([]CustomField, error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/custom_fields.json?%s",
			c.endpoint,
			c.getPaginationClause()),
		nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Redmine-API-Key", c.auth.Token)

	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r customFieldsResult
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
	return r.CustomFields, nil
}

package redmine

import (
	"encoding/json"
	"errors"
	"strings"
)

type timeEntryActivitiesResult struct {
	TimeEntryActivites []TimeEntryActivity `json:"time_entry_activities"`
}

type TimeEntryActivity struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	IsDefault bool   `json:"is_default"`
}

func (c *Client) TimeEntryActivities() ([]TimeEntryActivity, error) {
	res, err := c.Get(c.endpoint + "/enumerations/time_entry_activities.json?" + c.apiKeyParameter() + c.getPaginationClause())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r timeEntryActivitiesResult
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
	return r.TimeEntryActivites, nil
}

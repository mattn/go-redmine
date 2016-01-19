package redmine

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type issueRequest struct {
	Issue Issue `json:"issue"`
}

type issueResult struct {
	Issue Issue `json:"issue"`
}

type issuesResult struct {
	Issues []Issue `json:"issues"`
}

type Issue struct {
	Id          int     `json:"id"`
	Subject     string  `json:"subject"`
	Description string  `json:"description"`
	ProjectId   int     `json:"project_id"`
	Project     *IdName `json:"project"`
	Tracker     *IdName `json:"tracker"`
	StatusId    int     `json:"status_id"`
	Status      *IdName `json:"status"`
	Priority    *IdName `json:"priority"`
	Author      *IdName `json:"author"`
	AssignedTo  *IdName `json:"assigned_to"`
	Notes       string  `json:"notes"`
	StatusDate  string  `json:"status_date"`
	CreatedOn   string  `json:"created_on"`
	UpdatedOn   string  `json:"updated_on"`
	CustomFields []*CustomField `json:"custom_fields"`
}

type CustomField struct {
	Id      int     `json:"id"`
	Name	string  `json:"name"`
	Value	string  `json:"value"`
}

func (c *client) IssuesOf(projectId int) ([]Issue, error) {
	res, err := c.Get(c.endpoint + "/issues.json?project_id=" + strconv.Itoa(projectId) + "&key=" + c.apikey)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r issuesResult
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
	return r.Issues, nil
}

func (c *client) Issue(id int) (*Issue, error) {
	res, err := c.Get(c.endpoint + "/issues/" + strconv.Itoa(id) + ".json?key=" + c.apikey)
	if res.StatusCode == 404 {
		return nil, errors.New("Not Found")
	}
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r issueRequest
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
	return &r.Issue, nil
}

func (c *client) IssuesByQuery(query_id int) ([]Issue, error) {
	res, err := http.Get(c.endpoint + "/issues.json?query_id=" + strconv.Itoa(query_id) + "&key=" + c.apikey)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r issuesResult
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
	return r.Issues, nil
}

func (c *client) Issues() ([]Issue, error) {
	res, err := c.Get(c.endpoint + "/issues.json?key=" + c.apikey)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r issuesResult
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
	return r.Issues, nil
}

func (c *client) CreateIssue(issue Issue) (*Issue, error) {
	var ir issueRequest
	ir.Issue = issue
	s, err := json.Marshal(ir)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", c.endpoint+"/issues.json?key="+c.apikey, strings.NewReader(string(s)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r issueRequest
	if res.StatusCode != 201 {
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
	return &r.Issue, nil
}

func (c *client) UpdateIssue(issue Issue) error {
	var ir issueRequest
	ir.Issue = issue
	s, err := json.Marshal(ir)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", c.endpoint+"/issues/"+strconv.Itoa(issue.Id)+".json?key="+c.apikey, strings.NewReader(string(s)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.Do(req)
	if res.StatusCode == 404 {
		return errors.New("Not Found")
	}
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		decoder := json.NewDecoder(res.Body)
		var er errorsResult
		err = decoder.Decode(&er)
		if err == nil {
			err = errors.New(strings.Join(er.Errors, "\n"))
		}
	}
	if err != nil {
		return err
	}
	return err
}

func (c *client) DeleteIssue(id int) error {
	req, err := http.NewRequest("DELETE", c.endpoint+"/issues/"+strconv.Itoa(id)+".json?key="+c.apikey, strings.NewReader(""))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.Do(req)
	if res.StatusCode == 404 {
		return errors.New("Not Found")
	}
	if err != nil {
		return err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	if res.StatusCode != 200 {
		var er errorsResult
		err = decoder.Decode(&er)
		if err == nil {
			err = errors.New(strings.Join(er.Errors, "\n"))
		}
	}
	return err
}

func (issue *Issue) GetTitle() string {
	return issue.Tracker.Name + " #" + strconv.Itoa(issue.Id) + ": " + issue.Subject
}

package redmine

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type issueCategoriesResult struct {
	IssueCategories []IssueCategory `json:"categories"`
}

type issueCategoryResult struct {
	IssueCategory IssueCategory `json:"issue_category"`
}

type issueCategoryRequest struct {
	IssueCategory IssueCategory `json:"issue_category"`
}

type IssueCategory struct {
	Id      int    `json:"id"`
	Project IdName `json:"project"`
	Name    string `json:"name"`
}

func (c *client) IssueCategories(projectId int) ([]IssueCategory, error) {
	res, err := c.Get(c.endpoint + "/projects/" + strconv.Itoa(projectId) + "/issue_categories.json?key=" + c.apikey)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r issueCategoriesResult
	if res.StatusCode == 404 {
		return nil, errors.New("Not Found")
	}
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
	return r.IssueCategories, nil
}

func (c *client) IssueCategory(id int) (*IssueCategory, error) {
	res, err := c.Get(c.endpoint + "/issue_categories/" + strconv.Itoa(id) + ".json?key=" + c.apikey)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r issueCategoryResult
	if res.StatusCode == 404 {
		return nil, errors.New("Not Found")
	}
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
	return &r.IssueCategory, nil
}

func (c *client) CreateIssueCategory(issueCategory IssueCategory) (*IssueCategory, error) {
	var ir issueCategoryRequest
	ir.IssueCategory = issueCategory
	s, err := json.Marshal(ir)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", c.endpoint+"/issue_categories.json?key="+c.apikey, strings.NewReader(string(s)))
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
	var r issueCategoryResult
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
	return &r.IssueCategory, nil
}

func (c *client) UpdateIssueCategory(issueCategory IssueCategory) error {
	var ir issueCategoryRequest
	ir.IssueCategory = issueCategory
	s, err := json.Marshal(ir)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", c.endpoint+"/issue_categories/"+strconv.Itoa(issueCategory.Id)+".json?key="+c.apikey, strings.NewReader(string(s)))
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

func (c *client) DeleteIssueCategory(id int) error {
	req, err := http.NewRequest("DELETE", c.endpoint+"/issue_categories/"+strconv.Itoa(id)+".json?key="+c.apikey, strings.NewReader(""))
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

package redmine

import (
	"encoding/json"
	"errors"
	"fmt"
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
	Issues     []Issue `json:"issues"`
	TotalCount uint    `json:"total_count"`
	Offset     uint    `json:"offset"`
	Limit      uint    `json:"limit"`
}

type JournalDetails struct {
	Property string `json:"property"`
	Name     string `json:"name"`
	OldValue string `json:"old_value"`
	NewValue string `json:"new_value"`
}
type Journal struct {
	Id        int              `json:"id"`
	User      *IdName          `json:"user"`
	Notes     string           `json:"notes"`
	CreatedOn string           `json:"created_on"`
	Details   []JournalDetails `json:"details"`
}

type Issue struct {
	Id           int            `json:"id"`
	Subject      string         `json:"subject"`
	Description  string         `json:"description"`
	ProjectId    int            `json:"project_id"`
	Project      *IdName        `json:"project"`
	TrackerId    int            `json:"tracker_id"`
	Tracker      *IdName        `json:"tracker"`
	ParentId     int            `json:"parent_issue_id,omitempty"`
	Parent       *Id            `json:"parent"`
	StatusId     int            `json:"status_id"`
	Status       *IdName        `json:"status"`
	PriorityId   int            `json:"priority_id,omitempty"`
	Priority     *IdName        `json:"priority"`
	Author       *IdName        `json:"author"`
	FixedVersion *IdName        `json:"fixed_version"`
	AssignedTo   *IdName        `json:"assigned_to"`
	AssignedToId int            `json:"assigned_to_id"`
	Category     *IdName        `json:"category"`
	CategoryId   int            `json:"category_id"`
	Notes        string         `json:"notes"`
	StatusDate   string         `json:"status_date"`
	CreatedOn    string         `json:"created_on"`
	UpdatedOn    string         `json:"updated_on"`
	StartDate    string         `json:"start_date"`
	DueDate      string         `json:"due_date"`
	ClosedOn     string         `json:"closed_on"`
	CustomFields []*CustomField `json:"custom_fields,omitempty"`
	Uploads      []*Upload      `json:"uploads"`
	DoneRatio    float32        `json:"done_ratio"`
	EstimatedHours float32      `json:"estimated_hours"`
	Journals     []*Journal     `json:"journals"`
}

type IssueFilter struct {
	ProjectId    string
	SubprojectId string
	TrackerId    string
	StatusId     string
	AssignedToId string
	UpdatedOn    string
	ExtraFilters map[string]string
}

type CustomField struct {
	Id          int         `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Multiple    bool        `json:"multiple"`
	Value       interface{} `json:"value"`
}

func (c *Client) IssuesOf(projectId int) ([]Issue, error) {
	issues, err := getIssues(c, "/issues.json?project_id="+strconv.Itoa(projectId)+"&key="+c.apikey+c.getPaginationClause())

	if err != nil {
		return nil, err
	}

	return issues, nil
}

func (c *Client) Issue(id int) (*Issue, error) {
	return getOneIssue(c, id, nil)
}

func (c *Client) IssueWithArgs(id int, args map[string]string) (*Issue, error) {
	return getOneIssue(c, id, args)
}

func (c *Client) IssuesByQuery(queryId int) ([]Issue, error) {
	issues, err := getIssues(c, "/issues.json?query_id="+strconv.Itoa(queryId)+"&key="+c.apikey+c.getPaginationClause())

	if err != nil {
		return nil, err
	}
	return issues, nil
}

// IssuesByFilter filters issues applying the f criteria
func (c *Client) IssuesByFilter(f *IssueFilter) ([]Issue, error) {
	issues, err := getIssues(c, "/issues.json?key="+c.apikey+c.getPaginationClause()+getIssueFilterClause(f))
	if err != nil {
		return nil, err
	}
	return issues, nil
}

func (c *Client) Issues() ([]Issue, error) {
	issues, err := getIssues(c, "/issues.json?key="+c.apikey+c.getPaginationClause())

	if err != nil {
		return nil, err
	}

	return issues, nil
}

func (c *Client) CreateIssue(issue Issue) (*Issue, error) {
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

func (c *Client) UpdateIssue(issue Issue) error {
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
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return errors.New("Not Found")
	}
	if res.StatusCode/100 != 2 {
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

func (c *Client) DeleteIssue(id int) error {
	req, err := http.NewRequest("DELETE", c.endpoint+"/issues/"+strconv.Itoa(id)+".json?key="+c.apikey, strings.NewReader(""))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return errors.New("Not Found")
	}

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

// MarshalJSON marshals issue to JSON.
// This overrides the default MarshalJSON() to reset parent issue.
func (issue Issue) MarshalJSON() ([]byte, error) {
	type Issue2 Issue

	// To reset parent issue, set empty string to "parent_issue_id"
	var parentIssueID *string
	if issue.Parent == nil {
		// reset parent issue
		id := ""
		parentIssueID = &id
	} else if issue.ParentId > 0 {
		// set parent issue
		id := strconv.Itoa(issue.ParentId)
		parentIssueID = &id
	}

	return json.Marshal(&struct {
		Issue2
		ParentId *string `json:"parent_issue_id,omitempty"`
	}{
		Issue2:   Issue2(issue),
		ParentId: parentIssueID,
	})
}

func getIssueFilterClause(filter *IssueFilter) string {
	if filter == nil {
		return ""
	}
	clause := ""
	if filter.ProjectId != "" {
		clause = clause + fmt.Sprintf("&project_id=%v", filter.ProjectId)
	}
	if filter.SubprojectId != "" {
		clause = clause + fmt.Sprintf("&subproject_id=%v", filter.SubprojectId)
	}
	if filter.TrackerId != "" {
		clause = clause + fmt.Sprintf("&tracker_id=%v", filter.TrackerId)
	}
	if filter.StatusId != "" {
		clause = clause + fmt.Sprintf("&status_id=%v", filter.StatusId)
	}
	if filter.AssignedToId != "" {
		clause = clause + fmt.Sprintf("&assigned_to_id=%v", filter.AssignedToId)
	}
	if filter.UpdatedOn != "" {
		clause = clause + fmt.Sprintf("&updated_on=%v", filter.UpdatedOn)
	}

	if filter.ExtraFilters != nil {
		extraFilter := make([]string, 0)
		for key, value := range filter.ExtraFilters {
			extraFilter = append(extraFilter, fmt.Sprintf("%s=%s", key, value))
		}
		clause = clause + "&" + strings.Join(extraFilter[:], "&")
	}

	return clause
}

func mapConcat(m map[string]string, delimiter string) string {
	var args []string

	for k, v := range m {
		args = append(args, k+"="+v)
	}

	return strings.Join(args, delimiter)
}

func getOneIssue(c *Client, id int, args map[string]string) (*Issue, error) {
	url := c.endpoint + "/issues/" + strconv.Itoa(id) + ".json?key=" + c.apikey

	if args != nil {
		url += "&" + mapConcat(args, "&")
	}

	res, err := c.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return nil, errors.New("Not Found")
	}

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

func getIssue(c *Client, url string, offset int) (*issuesResult, error) {
	res, err := c.Get(c.endpoint + url + "&offset=" + strconv.Itoa(offset))

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

	return &r, nil
}

func getIssues(c *Client, url string) ([]Issue, error) {
	completed := false
	var issues []Issue

	for completed == false {
		r, err := getIssue(c, url, len(issues))

		if err != nil {
			return nil, err
		}

		if r.TotalCount == uint(len(issues)) {
			completed = true
		}

		issues = append(issues, r.Issues...)
	}

	return issues, nil
}

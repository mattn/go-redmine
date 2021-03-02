package redmine

import (
	"encoding/json"
	"errors"
	"fmt"
	errors2 "github.com/pkg/errors"
	"net/http"
	"strconv"
	"strings"
)

type projectRequest struct {
	Project Project `json:"project"`
}

type projectResult struct {
	Project Project `json:"project"`
}

type projectsResult struct {
	Projects []Project `json:"projects"`
}

// Project contains a Redmine API project object according Redmine 4.1 REST API.
//
// See also: https://www.redmine.org/projects/redmine/wiki/Rest_api
type Project struct {
	// Id uniquely identifies a project on technical level. This value will be generated on project creation and cannot
	// be changed. Id is mandatory for all project API calls except CreateProject()
	Id int `json:"id"`
	// ParentID may contain the Id of a parent project. If set, this project is then a child project of the parent project.
	// Projects can be unlimitedly nested.
	ParentID Id `json:"parent_id"`
	// Name contains a human readable project name.
	Name string `json:"name"`
	// Identifier used by the application for various things (eg. in URLs). It must be unique and cannot be composed of
	// only numbers. It must contain 1 to 100 characters of which only consist of lowercase latin characters, numbers,
	// hyphen (-) and underscore (_). Once the project is created, this identifier cannot be modified
	Identifier string `json:"identifier"`
	// Description contains a human readable project multiline description that appears on the project overview.
	Description string `json:"description"`
	// Homepage contains a URL to a project's website that appears on the project overview.
	Homepage string `json:"homepage"`
	// IsPublic controls who can view the project. If set to true the project can be viewed by all the users, including
	// those who are not members of the project. If set to false, only the project members have access to it, according to
	// their role.
	//
	// since Redmine 2.6.0
	IsPublic bool `json:"is_public"`
	// InheritMembers determines whether this project inherits members from a parent project. If set to true (and being a
	// nested project) all members from the parent project will apply also to this project.
	InheritMembers bool `json:"inherit_members"`
	// CreatedOn contains a timestamp of when the project was created.
	CreatedOn string `json:"created_on"`
	// UpdatedOn contains the timestamp of when the project was last updated.
	UpdatedOn string `json:"updated_on"`
}

// Project returns a single project without additional fields.
func (c *Client) Project(id int) (*Project, error) {
	req, err := c.authenticatedGet(c.endpoint + "/projects/" + strconv.Itoa(id) + ".json")
	if err != nil {
		return nil, errors2.Wrapf(err, "error while creating GET request for project %d ", id)
	}

	res, err := c.Do(req)
	if err != nil {
		return nil, errors2.Wrapf(err, "could not read project %d ", id)
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r projectResult
	if res.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("project (id: %d) was not found", id)
	}
	if !isHTTPStatusSuccessful(res.StatusCode, []int{http.StatusOK}) {
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
	return &r.Project, nil
}

func isHTTPStatusSuccessful(httpStatus int, acceptedStatuses []int) bool {
	for _, acceptedStatus := range acceptedStatuses {
		if httpStatus == acceptedStatus {
			return true
		}
	}

	return false
}

func (c *Client) Projects() ([]Project, error) {
	req, err := c.authenticatedGet(c.endpoint + "/projects.json")
	if err != nil {
		return nil, errors2.Wrap(err, "error while creating GET request for projects")
	}
	err = safelyAddQueryParameters(req, c.getPaginationClauseParams())
	if err != nil {
		return nil, errors2.Wrap(err, "error while adding pagination parameters to project request")
	}

	res, err := c.Do(req)
	if err != nil {
		return nil, errors2.Wrap(err, "could not read projects")
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r projectsResult
	if !isHTTPStatusSuccessful(res.StatusCode, []int{http.StatusOK}) {
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
	return r.Projects, nil
}

func (c *Client) CreateProject(project Project) (*Project, error) {
	var ir projectRequest
	ir.Project = project
	s, err := json.Marshal(ir)
	if err != nil {
		return nil, err
	}

	req, err := c.authenticatedRequest("POST", c.endpoint+"/projects.json", strings.NewReader(string(s)))
	if err != nil {
		return nil, errors2.Wrapf(err, "error while creating POST request for project %s ", project.Identifier)
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.Do(req)
	if err != nil {
		return nil, errors2.Wrapf(err, "could not create project %s ", project.Identifier)
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r projectRequest
	if !isHTTPStatusSuccessful(res.StatusCode, []int{http.StatusCreated}) {
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
	return &r.Project, nil
}

func (c *Client) UpdateProject(project Project) error {
	var ir projectRequest
	ir.Project = project
	s, err := json.Marshal(ir)
	if err != nil {
		return err
	}

	req, err := c.authenticatedRequest("PUT", c.endpoint+"/projects/"+strconv.Itoa(project.Id)+".json", strings.NewReader(string(s)))
	if err != nil {
		return errors2.Wrapf(err, "error while creating PUT request for project %d ", project.Id)
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.Do(req)
	if err != nil {
		return errors2.Wrapf(err, "could not update project %d ", project.Id)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return fmt.Errorf("could not update project (id: %d) because it was not found", project.Id)
	}
	if !isHTTPStatusSuccessful(res.StatusCode, []int{http.StatusOK, http.StatusNoContent}) {
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

func (c *Client) DeleteProject(id int) error {
	req, err := c.authenticatedRequest("DELETE", c.endpoint+"/projects/"+strconv.Itoa(id)+".json", strings.NewReader(""))
	if err != nil {
		return errors2.Wrapf(err, "error while creating DELETE for project %d ", id)
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.Do(req)
	if err != nil {
		return errors2.Wrapf(err, "could not delete project %d ", id)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return fmt.Errorf("could not delete project (id %d) because it was not found", id)
	}

	decoder := json.NewDecoder(res.Body)
	if !isHTTPStatusSuccessful(res.StatusCode, []int{http.StatusOK, http.StatusNoContent}) {
		var er errorsResult
		err = decoder.Decode(&er)
		if err == nil {
			err = errors.New(strings.Join(er.Errors, "\n"))
		}
	}
	return err
}

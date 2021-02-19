package redmine

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const (
	// since Redmine 2.6.0
	ProjectAdditionalFieldTrackers = "trackers"
	// since Redmine 2.6.0
	ProjectAdditionalFieldIssueCategories = "issue_categories"
	// since Redmine 2.6.0
	ProjectAdditionalFieldEnabledModules = "enabled_modules"
	// since Redmine 3.4.0
	ProjectAdditionalFieldTimeEntryActivities = "time_entry_activities"
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
	// TrackerIDs
	// since Redmine 2.6.0
	TrackerIDs []int `json:"tracker_ids"`
	// EnabledModuleNames
	// since Redmine 2.6.0
	EnabledModuleNames []string `json:"enabled_module_names"`
	// IssueCategories
	// since Redmine 2.6.0
	IssueCategories IssueCategoriesResult `json:"issue_categories"`
	// CustomFields.
	// the Redmine API description is unclear about this field: Docu mentions issue_custom_field_ids only for
	// project creation.
	CustomFields []int `json:"issue_custom_field_ids,omitempty"`
	// CreatedOn contains a timestamp of when the project was created.
	CreatedOn string `json:"created_on"`
	// UpdatedOn contains the timestamp of when the project was last updated.
	UpdatedOn string `json:"updated_on"`
}

// Project returns a single project without additional fields.
func (c *Client) Project(id int) (*Project, error) {
	return c.ProjectWithAdditionalFields(id)
}

// ProjectWithAdditionalFields returns a single project along with additional fields selected by the caller. The given
// additional fields can be nil, empty or a set of the currently supported additional project fields.
//
// Example to include trackers:
//  project, err := client.ProjectWithAdditionalFields(42, redmine, ProjectAdditionalFieldTrackers, ProjectAdditionalFieldEnabledModules)
func (c *Client) ProjectWithAdditionalFields(id int, additionalFields ...string) (*Project, error) {
	err := validateAdditionalFields(additionalFields...)
	if err != nil {
		return nil, err
	}
	additionalFieldsParameter := additionalFieldsParam(additionalFieldsParam())
	parameters := c.concatParameters(c.apiKeyParameter(), additionalFieldsParameter)

	res, err := c.Get(c.endpoint + "/projects/" + strconv.Itoa(id) + ".json?" + parameters)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r projectResult
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

// validateAdditionalFields checks for invalid field names. Repeated fields will be ignored because Redmine handles
// multiple instances of the same field without error.
func validateAdditionalFields(additionalFields ...string) error {
	for _, field := range additionalFields {
		if field == ProjectAdditionalFieldTrackers ||
			field == ProjectAdditionalFieldIssueCategories ||
			field == ProjectAdditionalFieldEnabledModules ||
			field == ProjectAdditionalFieldTimeEntryActivities {
			continue
		}
		return fmt.Errorf("unsupported additional project field %s found", field)
	}

	return nil
}

func additionalFieldsParam(additionalFields ...string) string {
	return strings.Join(additionalFields, ",")
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
	return c.ProjectsWithAdditionalFields()
}

func (c *Client) ProjectsWithAdditionalFields(additionalFields ...string) ([]Project, error) {
	err := validateAdditionalFields(additionalFields...)
	if err != nil {
		return nil, err
	}
	additionalFieldsParameter := additionalFieldsParam(additionalFieldsParam())

	parameters := c.concatParameters(c.apiKeyParameter(), additionalFieldsParameter, c.getPaginationClause())
	res, err := c.Get(c.endpoint + "/projects.json?" + parameters)
	if err != nil {
		return nil, err
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

	parameters := c.concatParameters(c.apiKeyParameter())
	req, err := http.NewRequest("POST", c.endpoint+"/projects.json?"+parameters, strings.NewReader(string(s)))
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

	parameters := c.concatParameters(c.apiKeyParameter())
	req, err := http.NewRequest("PUT", c.endpoint+"/projects/"+strconv.Itoa(project.Id)+".json?"+parameters, strings.NewReader(string(s)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return fmt.Errorf("could not update project (id %d) because it was not found", project.Id)
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
	parameters := c.concatParameters(c.apiKeyParameter())
	req, err := http.NewRequest("DELETE", c.endpoint+"/projects/"+strconv.Itoa(id)+".json?"+parameters, strings.NewReader(""))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return fmt.Errorf("could not delete project (id %d) because it was not found", id)
	}

	decoder := json.NewDecoder(res.Body)
	if !isHTTPStatusSuccessful(res.StatusCode, []int{http.StatusOK}) {
		var er errorsResult
		err = decoder.Decode(&er)
		if err == nil {
			err = errors.New(strings.Join(er.Errors, "\n"))
		}
	}
	return err
}

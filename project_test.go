package redmine

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_Project(t *testing.T) {
	t.Run("should parse general project fields, and ignore module names and trackers from http response", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprintln(w, `{
  "project": {
    "id": 1,
    "name": "example project",
    "identifier": "exampleproject",
    "description": "This is an example project.",
    "homepage": "http://github.com/cloudogu/go-redmine",
    "status": 1,
    "is_public": true,
    "inherit_members": true,
    "trackers": [
      {
        "id": 1,
        "name": "Bug"
      },
      {
        "id": 2,
        "name": "Feature"
      }
    ],
    "enabled_modules": [
      {
        "id": 71,
        "name": "issue_tracking"
      },
			{
        "id": 73,
        "name": "wiki"
      }
    ],
    "created_on": "2021-02-19T16:51:03Z",
    "updated_on": "2021-02-19T16:51:25Z"
  }
}`)
		}))
		defer ts.Close()

		sut, err := NewClient(ts.URL, APIAuth{AuthType: AuthTypeTokenQueryParam, Token: "apiKey"})

		actualProject, err := sut.Project(1)

		require.NoError(t, err)
		require.NotEmpty(t, actualProject)
		expectedProject := &Project{
			Id:             1,
			ParentID:       Id{},
			Name:           "example project",
			Identifier:     "exampleproject",
			Description:    "This is an example project.",
			Homepage:       "http://github.com/cloudogu/go-redmine",
			IsPublic:       true,
			InheritMembers: true,
			CreatedOn:      "2021-02-19T16:51:03Z",
			UpdatedOn:      "2021-02-19T16:51:25Z",
		}
		assert.Equal(t, expectedProject, actualProject)
	})
}

package redmine

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_validateAdditionalFields(t *testing.T) {
	type args struct {
		fields []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"should pass for single valid field", args{fields: []string{"trackers"}}, false},
		{"should pass for multiple valid fields", args{fields: []string{"enabled_modules", "trackers"}}, false},
		{"should pass for all fields", args{fields: []string{
			"trackers", "issue_categories", "enabled_modules", "time_entry_activities"}}, false},
		{"should pass for repeated fields", args{fields: []string{
			"trackers", "issue_categories", "enabled_modules", "time_entry_activities",
			"trackers", "issue_categories", "enabled_modules", "time_entry_activities"}}, false},
		{"should fail for single unsupported field", args{fields: []string{"invalid value"}}, true},
		{"should fail for single unsupported field among valid fields", args{fields: []string{"enabled_modules", "invalidvalue"}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateAdditionalFields(tt.args.fields...); (err != nil) != tt.wantErr {
				t.Errorf("validateAdditionalFields() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	t.Run("should pass for non-existing parameter", func(t *testing.T) {
		if err := validateAdditionalFields(); (err != nil) != false {
			t.Errorf("expected error for empty string")
		}
	})
	t.Run("should fail for empty string", func(t *testing.T) {
		if err := validateAdditionalFields(""); (err != nil) != true {
			t.Errorf("expected error for empty string")
		}
	})
}

func Test_additionalFieldsParam(t *testing.T) {
	type args struct {
		additionalFields []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"should create empty string for nil", args{nil}, ""},
		{"should create empty string for empty slice", args{[]string{}}, ""},
		{"should create parameter without comma for single field", args{[]string{"enabled_modules"}}, "enabled_modules"},
		{"should create comma delimited list without trailing comma for multiple fields", args{[]string{
			"enabled_modules", "enabled_modules", "trackers"}}, "enabled_modules,enabled_modules,trackers"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := additionalFieldsParam(tt.args.additionalFields...); got != tt.want {
				t.Errorf("additionalFieldsParam() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_ProjectWithAdditionalFields(t *testing.T) {
	t.Run("should parse general project fields, module names, and trackers from http response", func(t *testing.T) {
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

		sut := NewClient(ts.URL, "apiKey")

		actualProject, err := sut.ProjectWithAdditionalFields(1, "enabled_modules", "trackers")

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
			Trackers: []Tracker{
				{ID: 1, Name: "Bug"},
				{ID: 2, Name: "Feature"},
			},
			EnabledModules: []EnabledModul{
				{ID: 71, Name: "issue_tracking"},
				{ID: 73, Name: "wiki"},
			},
			CreatedOn: "2021-02-19T16:51:03Z",
			UpdatedOn: "2021-02-19T16:51:25Z",
		}
		assert.Equal(t, expectedProject, actualProject)
	})
}

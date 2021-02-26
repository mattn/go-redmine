package redmine

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_getOneIssue(t *testing.T) {
	t.Run("should parse simple issue JSON without additional arguments", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprintln(w, `{
  "issue": {
    "id": 1,
    "project": {
      "id": 1,
      "name": "example project1"
    },
    "tracker": {
      "id": 1,
      "name": "Bug"
    },
    "status": {
      "id": 1,
      "name": "New"
    },
    "priority": {
      "id": 2,
      "name": "Normal"
    },
    "author": {
      "id": 1,
      "name": "Redmine Admin"
    },
    "subject": "Something should be done",
    "description": "In this ticket an **important task** should be done1!\r\n\r\nGo ahead!\r\n\r\n`+"```bash\\r\\necho -n $PATH\\r\\n```"+`",
    "start_date": null,
    "due_date": null,
    "done_ratio": 0,
    "is_private": false,
    "estimated_hours": null,
    "total_estimated_hours": null,
    "spent_hours": 0,
    "total_spent_hours": 0,
    "created_on": "2021-02-23T14:20:48Z",
    "updated_on": "2021-02-23T14:39:02Z",
    "closed_on": null
  }
}`)
		}))
		defer ts.Close()

		sut, _ := NewClient(ts.URL, APIAuth{AuthType: AuthTypeTokenQueryParam, Token: "apiKey"})

		actual, err := getOneIssue(sut, 1, nil)

		require.NoError(t, err)
		assert.Equal(t, 1, actual.Id)
		assert.Equal(t, "Something should be done", actual.Subject)
		assert.Equal(t, "In this ticket an **important task** should be done1!\r\n\r\nGo ahead!\r\n\r\n"+"```bash\r\necho -n $PATH\r\n```", actual.Description)
		assert.Equal(t, 0, actual.ProjectId)
		assert.Equal(t, 0, actual.TrackerId)
		assert.Equal(t, 0, actual.ParentId)
		assert.Equal(t, 0, actual.StatusId)
		assert.Equal(t, 0, actual.PriorityId)
		assert.Equal(t, "2021-02-23T14:20:48Z", actual.CreatedOn)
		assert.Equal(t, "2021-02-23T14:39:02Z", actual.UpdatedOn)
		assert.Equal(t, "", actual.StartDate)
		assert.Equal(t, "", actual.DueDate)
		assert.Equal(t, "", actual.ClosedOn)

		expectedProject := IdName{Id: 1, Name: "example project1"}
		assert.Equal(t, expectedProject, *actual.Project)
		expectedTracker := IdName{Id: 1, Name: "Bug"}
		assert.Equal(t, expectedTracker, *actual.Tracker)
		assert.Nil(t, actual.Parent)
		expectedStatus := IdName{Id: 1, Name: "New"}
		assert.Equal(t, expectedStatus, *actual.Status)
		expectedPriority := IdName{Id: 2, Name: "Normal"}
		assert.Equal(t, expectedPriority, *actual.Priority)
		expectedAuthor := IdName{Id: 1, Name: "Redmine Admin"}
		assert.Equal(t, expectedAuthor, *actual.Author)
	})
}

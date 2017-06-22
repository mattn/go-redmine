package redmine

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type wikiPagesResult struct {
	WikiPages []WikiPage `json:"wiki_pages"`
}

type wikiPageResult struct {
	WikiPage WikiPage `json:"wiki_page"`
}

type wikiPageRequest struct {
	WikiPage WikiPage `json:"wiki_page"`
}

type WikiPage struct {
	Title     string      `json:"title"`
	Parent    *Parent     `json:"parent,omitempty"`
	Text      string      `json:"text"`
	Version   interface{} `json:"version,omitempty"`
	Author    *IdName     `json:"author,omitempty"`
	Comments  string      `json:"comments"`
	CreatedOn string      `json:"created_on,omitempty"`
	UpdatedOn string      `json:"updated_on,omitempty"`
	ParentID  int         `json:"parent_id"`
}

type Parent struct {
	Title string `json:"title"`
}

// WikiPages fetches a list of all wiki pages of the given project.
// The Text field of the listed pages is not fetch by this command and is thus empty.
func (c *Client) WikiPages(projectId int) ([]WikiPage, error) {
	res, err := c.Get(c.endpoint + "/projects/" + strconv.Itoa(projectId) + "/wiki/index.json?key=" + c.apikey + c.getPaginationClause())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r wikiPagesResult
	if res.StatusCode == 404 {
		return nil, errors.New("Not Found")
	}
	if res.StatusCode != 200 {
		var er errorsResult
		if err = decoder.Decode(&er); err != nil {
			return nil, err
		}
		return nil, errors.New(strings.Join(er.Errors, "\n"))
	} else {
		if err = decoder.Decode(&r); err != nil {
			return nil, err
		}
	}
	return r.WikiPages, nil
}

// WikiPage fetches the wiki page with the given title.
func (c *Client) WikiPage(projectId int, title string) (*WikiPage, error) {
	return c.getWikiPage(projectId, title)
}

// WikiPageAtVersion fetches the wiki page with the given title at the given version.
func (c *Client) WikiPageAtVersion(projectId int, title string, version string) (*WikiPage, error) {
	return c.getWikiPage(projectId, title+"/"+version)
}

func (c *Client) getWikiPage(projectId int, resource string) (*WikiPage, error) {
	res, err := c.Get(c.endpoint + "/projects/" + strconv.Itoa(projectId) + "/wiki/" + resource + ".json?key=" + c.apikey)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r wikiPageResult
	if res.StatusCode == 404 {
		return nil, errors.New("Not Found")
	}
	if res.StatusCode != 200 {
		var er errorsResult
		if err = decoder.Decode(&er); err != nil {
			return nil, err
		}
		return nil, errors.New(strings.Join(er.Errors, "\n"))
	} else {
		if err = decoder.Decode(&r); err != nil {
			return nil, err
		}
	}
	return &r.WikiPage, nil
}

// CreateWikiPage creates wiki page.
func (c *Client) CreateWikiPage(projectId int, wikiPage WikiPage) (*WikiPage, error) {
	var wpr wikiPageRequest
	wpr.WikiPage = wikiPage
	s, err := json.Marshal(wpr)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PUT", c.endpoint+"/projects/"+strconv.Itoa(projectId)+"/wiki/"+wikiPage.Title+".json?key="+c.apikey, strings.NewReader(string(s)))
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
	var r wikiPageResult
	if res.StatusCode != 201 {
		var er errorsResult
		if err = decoder.Decode(&er); err != nil {
			return nil, err
		}
		return nil, errors.New(strings.Join(er.Errors, "\n"))
	} else {
		if err := decoder.Decode(&r); err != nil {
			return nil, err
		}
	}
	return &r.WikiPage, nil
}

// UpdateWikiPage updates the wiki page given by the Title field of wikiPage.
func (c *Client) UpdateWikiPage(projectId int, wikiPage WikiPage) error {
	var wpr wikiPageRequest
	wpr.WikiPage = wikiPage
	s, err := json.Marshal(wpr)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", c.endpoint+"/projects/"+strconv.Itoa(projectId)+"/wiki/"+wikiPage.Title+".json?key="+c.apikey, strings.NewReader(string(s)))
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

	if res.StatusCode != 200 {
		decoder := json.NewDecoder(res.Body)
		var er errorsResult
		if err := decoder.Decode(&er); err != nil {
			return err
		}
		return errors.New(strings.Join(er.Errors, "\n"))
	}
	return nil
}

// DeleteWikiPage deletes the wiki page given by its title irreversibly.
func (c *Client) DeleteWikiPage(projectId int, title string) error {
	req, err := http.NewRequest("DELETE", c.endpoint+"/projects/"+strconv.Itoa(projectId)+"/wiki/"+title+".json?key="+c.apikey, strings.NewReader(""))
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
		if err := decoder.Decode(&er); err != nil {
			return err
		}
		return errors.New(strings.Join(er.Errors, "\n"))
	}
	return nil
}

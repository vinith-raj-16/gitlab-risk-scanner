package gitlab

import (
	"fmt"
)

type Project struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Path              string `json:"path"`
	WebURL            string `json:"web_url"`
	DefaultBranch     string `json:"default_branch"`
	Visibility        string `json:"visibility"`
	Archived          bool   `json:"archived"`
	Namespace         Namespace `json:"namespace"`
}

type Namespace struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
	Kind string `json:"kind"`
}

func (c *Client) GetUserProjects(username string) ([]Project, error) {
	endpoint := fmt.Sprintf(
		"/users/%s/projects?visibility=public&per_page=100",
		encode(username),
	)

	body, err := c.doRequest("GET", endpoint)
	if err != nil {
		return nil, err
	}

	var projects []Project

	if err := decodeJSON(body, &projects); err != nil {
		return nil, err
	}

	return projects, nil
}

func (c *Client) GetGroupProjects(group string) ([]Project, error) {
	endpoint := fmt.Sprintf(
		"/groups/%s/projects?visibility=public&include_subgroups=true&per_page=100",
		encode(group),
	)

	body, err := c.doRequest("GET", endpoint)
	if err != nil {
		return nil, err
	}

	var projects []Project

	if err := decodeJSON(body, &projects); err != nil {
		return nil, err
	}

	return projects, nil
}

func (c *Client) GetProject(fullPath string) (Project, error) {
	endpoint := fmt.Sprintf(
		"/projects/%s",
		encode(fullPath),
	)

	body, err := c.doRequest("GET", endpoint)
	if err != nil {
		return Project{}, err
	}

	var project Project

	if err := decodeJSON(body, &project); err != nil {
		return Project{}, err
	}

	return project, nil
}

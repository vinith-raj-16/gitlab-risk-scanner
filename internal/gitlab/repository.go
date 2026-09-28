package gitlab

import (
	"fmt"
	"net/url"
)

type TreeItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Path string `json:"path"`
	Mode string `json:"mode"`
}

func (c *Client) GetRepositoryTree(projectID int) ([]TreeItem, error) {
	endpoint := fmt.Sprintf(
		"/projects/%d/repository/tree?recursive=true&per_page=100",
		projectID,
	)

	body, err := c.doRequest("GET", endpoint)
	if err != nil {
		return nil, err
	}

	var items []TreeItem

	if err := decodeJSON(body, &items); err != nil {
		return nil, err
	}

	return items, nil
}

func (c *Client) GetFile(projectID int, filePath, branch string) (string, error) {
	endpoint := fmt.Sprintf(
		"/projects/%d/repository/files/%s?ref=%s",
		projectID,
		url.PathEscape(filePath),
		url.QueryEscape(branch),
	)

	body, err := c.doRequest("GET", endpoint)
	if err != nil {
		return "", err
	}

	var file struct {
		Content  string `json:"content"`
		Encoding string `json:"encoding"`
	}

	if err := decodeJSON(body, &file); err != nil {
		return "", err
	}

	// GitLab normally returns Base64 encoded content.
	if file.Encoding == "base64" {
		return decodeBase64(file.Content)
	}

	return file.Content, nil
}

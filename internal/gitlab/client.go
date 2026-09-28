package gitlab

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const baseURL = "https://gitlab.com/api/v4"

type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

func NewClient(token string) *Client {
	return &Client{
		BaseURL:    baseURL,
		Token:      token,
		HTTPClient: &http.Client{},
	}
}

func (c *Client) doRequest(method, endpoint string) ([]byte, error) {
	req, err := http.NewRequest(method, c.BaseURL+endpoint, nil)
	if err != nil {
		return nil, err
	}

	if c.Token != "" {
		req.Header.Set("PRIVATE-TOKEN", c.Token)
	}


	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf(
			"GitLab API returned status %d: %s",
			resp.StatusCode,
			string(body),
		)
	}

	return body, nil
}

func encode(value string) string {
	return url.PathEscape(value)
}

func decodeJSON(data []byte, target any) error {
	return json.Unmarshal(data, target)
}

func decodeBase64(value string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

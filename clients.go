package clients

import (
	"net/http"
	"net/url"
)

type BaseHTTPClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func (c *BaseHTTPClient) ResolvePath(path string) (string, error) {
	baseURL, err := url.Parse(c.BaseURL)
	if err != nil {
		return "", err
	}
	rel := baseURL.ResolveReference(&url.URL{Path: path})
	return rel.String(), err
}

func (c *BaseHTTPClient) DoRequest(req *http.Request) (*http.Response, error) {

	if c.HTTPClient == nil {
		c.HTTPClient = &http.Client{}
	}
	return c.HTTPClient.Do(req)
}

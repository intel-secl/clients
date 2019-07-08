package clients

import "net/url"

func ResolvePath(baseURLIn, path string) (string, error) {
	baseURL, err := url.Parse(baseURLIn)
	if err != nil {
		return "", err
	}
	rel := baseURL.ResolveReference(&url.URL{Path: path})
	return rel.String(), err
}

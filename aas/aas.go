package aas

import (
	"crypto/tls"
	"net/http"
	"net/url"
)

func httpClientNoTLS() *http.Client {

	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}

func resolvePath(baseURLIn, path string) (string, error) {
	baseURL, err := url.Parse(baseURLIn)
	if err != nil {
		return "", err
	}
	rel := baseURL.ResolveReference(&url.URL{Path: path})
	return rel.String(), err
}

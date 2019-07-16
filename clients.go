package clients

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
)

type HTTPClientErr struct {
	ErrMessage string
	ErrCode    int
}

func (ucErr *HTTPClientErr) Error() string {
	return fmt.Sprintf("%s: http status %d", ucErr.ErrMessage, ucErr.ErrCode)
}

func HTTPClient() *http.Client {
	return &http.Client{}
}

func HTTPClientTLSNoVerify() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}

func ResolvePath(baseURLIn, path string) (string, error) {
	baseURL, err := url.Parse(baseURLIn)
	if err != nil {
		return "", err
	}
	rel := baseURL.ResolveReference(&url.URL{Path: path})
	return rel.String(), err
}

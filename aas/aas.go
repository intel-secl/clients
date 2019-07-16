package aas

import (
	"crypto/tls"
	"net/http"
	"net/url"
)

type HTTPClientErr struct {
	ErrMessage string
	ErrCode    int
}

func (ucErr *ClientErr) Error() string {
	return fmt.Sprintf("%s: http status %d", ucErr.ErrMessage, ucErr.ErrCode)
}

func httpClient()*http.Client {
	return &http.Client{}
}

func resolvePath(baseURLIn, path string) (string, error) {
	baseURL, err := url.Parse(baseURLIn)
	if err != nil {
		return "", err
	}
	rel := baseURL.ResolveReference(&url.URL{Path: path})
	return rel.String(), err
}

package clients

import (
	"crypto/tls"
	"fmt"
	"net/http"
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

func ResolvePath(baseURL, path string) string {
	return baseURL + "/" + path
}

package clients

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

func HTTPClientWithCADir(caDir string) (*http.Client, error) {
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}
	var err error
	var allCA []string
	err = filepath.Walk(caDir, func(path string, info os.FileInfo, err error) error {
		allCA = append(allCA, path)
		return nil
	})
	if err != nil {
		return nil, err
	}
	for _, caFile := range allCA {
		ca, err := ioutil.ReadFile(caFile)
		if err == nil {
			rootCAs.AppendCertsFromPEM(ca)
		}
	}
	config := &tls.Config{
		InsecureSkipVerify: false,
		RootCAs:            rootCAs,
	}
	tr := &http.Transport{TLSClientConfig: config}

	return &http.Client{Transport: tr}, nil
}

func ResolvePath(baseURL, path string) string {
	if baseURL == "" ||
		path == "" {
		return ""
	}
	if strings.HasSuffix(baseURL, "/") {
                return baseURL + path
        }
	return baseURL +"/" + path
}

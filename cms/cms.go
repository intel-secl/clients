package cms

import (
	"bytes"
	"errors"
	"intel/isecl/lib/clients"
	"net/http"
)

type CMSClient struct {
	BaseUrl     string
	JWTToken    []byte
	httpClientP *http.Client
}

var (
	ErrFailToGetRootCA = errors.New("Failed to retrieve root CA")
	ErrSignCSRFailed   = errors.New("Failed to sign certificate with CMS")
)

func (c *CMSClient) httpClient() *http.Client {
	if c.httpClientP == nil {
		c.httpClientP = &http.Client{}
	}
	return c.httpClientP
}

func (c *CMSClient) GetRootCA() (string, error) {

	url, err := clients.ResolvePath(c.BaseUrl, "cms/v1/ca-certificates")
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Accept", "application/x-pem-file")
	rsp, err := c.httpClient().Do(req)
	if err != nil {
		return "", err
	}
	if rsp.StatusCode != http.StatusOK {
		return "", ErrFailToGetRootCA
	}
	resBuf := new(bytes.Buffer)
	resBuf.ReadFrom(rsp.Body)
	resStr := resBuf.String()
	return resStr, nil
}

func (c *CMSClient) PostCSR(csr []byte) (string, error) {

	url, err := clients.ResolvePath(c.BaseUrl, "cms/v1/certificates")
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(csr))

	req.Header.Set("Accept", "application/x-pem-file")
	req.Header.Set("Content-Type", "application/x-pem-file")

	req.Header.Add("Authorization", "Bearer "+string(c.JWTToken))
	rsp, err := c.httpClient().Do(req)
	if err != nil {
		return "", err
	}
	if rsp.StatusCode != http.StatusOK {
		return "", ErrSignCSRFailed
	}
	resBuf := new(bytes.Buffer)
	resBuf.ReadFrom(rsp.Body)
	resStr := resBuf.String()
	return resStr, nil
}

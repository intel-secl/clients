package cms

import (
	"bytes"
	"crypto/tls"
	"errors"
	"intel/isecl/lib/clients"
	"net/http"
)

type Client struct {
	BaseURL    string
	JWTToken   []byte
	HTTPClient *http.Client
}

var (
	ErrFailToGetRootCA = errors.New("Failed to retrieve root CA")
	ErrSignCSRFailed   = errors.New("Failed to sign certificate with CMS")
)

func (c *Client) httpClient() *http.Client {
	if c.HTTPClient == nil {
		tlsConfig := tls.Config{}
		tlsConfig.InsecureSkipVerify = true
		transport := http.Transport{
			TLSClientConfig: &tlsConfig,
		}
		c.HTTPClient = &http.Client{Transport: &transport}
	}
	return c.HTTPClient
}

func (c *Client) GetRootCA() (string, error) {

	url, err := clients.ResolvePath(c.BaseURL, "cms/v1/ca-certificates")
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

func (c *Client) PostCSR(csr []byte) (string, error) {

	url, err := clients.ResolvePath(c.BaseURL, "cms/v1/certificates")
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

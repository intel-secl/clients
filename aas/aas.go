package aas

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	
	sharedTypes "intel/isecl/lib/common/types/aas"
)

// Client for consuming the AAService API
type Client struct {
	BaseURL string
	// Username used to authenticate with the KMS. Username is only used for obtaining an authorization token, which is automatically used for requests.
	Username string
	// Password to supply for the Username
	Password string
	// Bearer Token
	JWTToken []byte
	// A reference to the underlying http Client.
	// If the value is nil, a default client will be created and used.
	HTTPClient *http.Client
}

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

func (c *Client) resolvePath(path string) (string, error) {
	baseURL, err := url.Parse(c.BaseURL)
	if err != nil {
		return "", err
	}
	rel := baseURL.ResolveReference(&url.URL{Path: path})
	return rel.String(), err
}

func (c *Client) dispatchRequest(req *http.Request) (*http.Response, error) {
	if c.Username != "" && c.Password != "" {
		req.SetBasicAuth(c.Username, c.Password)
	}
	if c.JWTToken != nil {
		req.Header.Add("Authorization", "Bearer "+string(c.JWTToken))
	}
	return c.httpClient().Do(req)
}

func (c *Client) GetJwtSigningCert() ([]byte, error) {
	jwtCertUrl, err := c.resolvePath("noauth/jwtCert")
	if err != nil {
		return nil, err
	}
	req, _ := http.NewRequest(http.MethodGet, jwtCertUrl, nil)
	req.Header.Set("Accept", "application/x-pem-file")
	rsp, err := c.dispatchRequest(req)
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to retrieve JWT signing certificate: HTTP Code: %d", rsp.StatusCode)
	}

	return ioutil.ReadAll(rsp.Body)
}

func (c *Client) GetJwtToken() ([]byte, error) {
	jwtUrl, err := c.resolvePath("token")
	if err != nil {
		return nil, err
	}
	req, _ := http.NewRequest(http.MethodGet, jwtUrl, nil)
	req.Header.Set("Accept", "application/jwt")
	rsp, err := c.dispatchRequest(req)
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to retrieve JWT (json web token): HTTP Code: %d", rsp.StatusCode)
	}

	return ioutil.ReadAll(rsp.Body)
}

func (c *Client) CreateUser(u sharedTypes.UserCreate) (*sharedTypes.UserCreateResponse, error) {
	userURL, err := c.resolvePath("users")
	if err != nil {
		return nil, err
	}

	payload, err := json.Marshal(&u)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, userURL, bytes.NewBuffer(payload))
	if err != nil {
		return nil, errors.New("Error creating request")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	rsp, err := c.dispatchRequest(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to create dispatch request: %v", err)
	}

	if rsp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Failed to create user: HTTP Code: %d", rsp.StatusCode)
	}
	var userCreateResponse sharedTypes.UserCreateResponse
	err = json.NewDecoder(rsp.Body).Decode(&userCreateResponse)
	if err != nil {
		return nil, err
	}
	return &userCreateResponse, nil
}

// CreateRole is used to create role in aas
func (c *Client) CreateRole(r sharedTypes.RoleCreate) (*sharedTypes.RoleCreateResponse, error) {

	roleURL, _ := c.resolvePath("roles")

	payload, err := json.Marshal(&r)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, roleURL, bytes.NewBuffer(payload))
	if err != nil {
		return nil, errors.New("Error creating request")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	rsp, err := c.dispatchRequest(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to create dispatch request: %v", err)
	}
	if rsp.StatusCode != http.StatusCreated {
		msg, _ := ioutil.ReadAll(rsp.Body)
		return nil, fmt.Errorf("Failed to create role: HTTP Code: %d: Response Message: %s", rsp.StatusCode, msg)
	}
	var roleCreateResponse sharedTypes.RoleCreateResponse
	err = json.NewDecoder(rsp.Body).Decode(&roleCreateResponse)
	if err != nil {
		return nil, err
	}
	return &roleCreateResponse, nil
}

// CreateRole is used to create role in aas
func (c *Client) AddRoleToUser(userID string, r sharedTypes.UserRoleCreate) error {
	reqURL, err := c.resolvePath("users" + userID + "roles")
	if err != nil {
		return err
	}

	payload, err := json.Marshal(&r)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, reqURL, bytes.NewBuffer(payload))
	if err != nil {
		return errors.New("Error creating request")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	rsp, err := c.dispatchRequest(req)
	if rsp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to add role user: HTTP Code: %d", rsp.StatusCode)
	}

	if err != nil {
		return err
	}
	return nil
}

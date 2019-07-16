package aas

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"intel/isecl/lib/clients"
	types "intel/isecl/lib/common/types/aas"
)

type JWTClientErr struct {
	ErrMessage string
	ErrInfo    string
}

func (ucErr *JWTClientErr) Error() string {
	return fmt.Sprintf("%s: %s", ucErr.ErrMessage, ucErr.ErrInfo)
}

var (
	ErrHTTPGetJWTCert = &clients.HTTPClientErr{
		ErrMessage: "Failed to retrieve JWT signing certificate",
		ErrCode:    0,
	}
	ErrHTTPFetchJWTToken = &clients.HTTPClientErr{
		ErrMessage: "Failed to retrieve JWT token from aas",
		ErrCode:    0,
	}
	ErrUserNotFound = &JWTClientErr{
		ErrMessage: "User name not registered",
		ErrInfo:    "",
	}
	ErrJWTNotYetFetched = &JWTClientErr{
		ErrMessage: "User token not yet fetched",
		ErrInfo:    "",
	}
)

type JWTClient struct {
	BaseURL        string
	NoTLSKeyVerify bool

	httpClientP *http.Client
	users       map[string]*types.UserCred
	tokens      map[string][]byte
}

func (c *JWTClient) GetJWTSigningCert() ([]byte, error) {

	jwtCertUrl, err := clients.ResolvePath(c.BaseURL, "noauth/jwtCert")
	if err != nil {
		return nil, err
	}
	req, _ := http.NewRequest(http.MethodGet, jwtCertUrl, nil)
	req.Header.Set("Accept", "application/x-pem-file")

	if c.httpClientP == nil {
		if c.NoTLSKeyVerify {
			c.httpClientP = clients.HTTPClientTLSNoVerify()
		} else {
			c.httpClientP = clients.HTTPClient()
		}
	}
	rsp, err := c.httpClientP.Do(req)
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode != http.StatusOK {
		ErrHTTPGetJWTCert.ErrCode = rsp.StatusCode
		return nil, ErrHTTPGetJWTCert
	}
	return ioutil.ReadAll(rsp.Body)
}

func (c *JWTClient) AddUser(username, password string) {
	c.users[username] = &types.UserCred{
		UserName: username,
		Password: password,
	}
}

func (c *JWTClient) GetUserToken(username string) ([]byte, error) {

	if _, ok := c.users[username]; !ok {
		ErrUserNotFound.ErrInfo = username
		return nil, ErrUserNotFound
	}
	token, ok := c.tokens[username]
	if ok {
		return token, nil
	}
	ErrJWTNotYetFetched.ErrInfo = username
	return nil, ErrJWTNotYetFetched
}

func (c *JWTClient) FetchAllTokens() error {

	for user, userCred := range c.users {
		token, err := c.fetchToken(userCred)
		if err != nil {
			return err
		}
		c.tokens[user] = token
	}
	return nil
}

func (c *JWTClient) FetchTokenForUser(username string) ([]byte, error) {

	userCred, ok := c.users[username]
	if !ok {
		return nil, ErrUserNotFound
	}
	token, err := c.fetchToken(userCred)
	if err != nil {
		return nil, err
	}
	c.tokens[username] = token
	return token, nil
}

func (c *JWTClient) fetchToken(userCred *types.UserCred) ([]byte, error) {

	jwtUrl, err := clients.ResolvePath(c.BaseURL, "token")
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(userCred)
	if err != nil {
		return nil, err
	}
	req, _ := http.NewRequest("POST", jwtUrl, buf)
	req.Header.Set("Accept", "application/jwt")

	if c.httpClientP == nil {
		if c.NoTLSKeyVerify {
			c.httpClientP = clients.HTTPClientTLSNoVerify()
		} else {
			c.httpClientP = clients.HTTPClient()
		}
	}
	rsp, err := c.httpClientP.Do(req)
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode != http.StatusOK {
		ErrHTTPFetchJWTToken.ErrCode = rsp.StatusCode
		return nil, ErrHTTPFetchJWTToken
	}
	jwtToken, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	return jwtToken, nil
}

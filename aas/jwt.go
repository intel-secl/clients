package aas

import (
	"bytes"
	"errors"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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
	ErrHTTPGetJWTCert  = &HTTPClientErr{"Failed to retrieve JWT signing certificate", 0}
	ErrHTTPFetchJWTToken = &HTTPClientErr{"Failed to retrieve JWT token from aas", 0}

	ErrUserNotFound  = &JWTClientErr{"User name not registered", ""}
	ErrJWTNotYetFetched  = &JWTClientErr{"User token not yet fetched", ""}
)

type JWTClient struct {
	BaseURL  string

	httpClientP *http.Client
	users map[string]*types.UserCred
	tokens map[string][]byte
}

func (c *JWTClient) GetJWTSigningCert() ([]byte, error) {

	jwtCertUrl, err := resolvePath(c.BaseURL, "noauth/jwtCert")
	if err != nil {
		return nil, err
	}
	req, _ := http.NewRequest(http.MethodGet, jwtCertUrl, nil)
	req.Header.Set("Accept", "application/x-pem-file")

	if c.httpClientP == nil {
		c.httpClientP = httpClient()
	}
	rsp, err := c.httpClientP.Do(req)
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode != http.StatusOK {
		ErrFailedToGetJWTCert.ErrCode = rsp.StatusCode
		return nil, ErrFailedToGetJWTCert
	}
	return ioutil.ReadAll(rsp.Body)
}

func (c *JWTClient) AddUser(username, password string) {
	users[username] = &types.UserCred{
		UserName: username,
		Password: password,
	}
}

func (c *JWTClient) GetUserToken(username string) ([]byte, error) {

	if _, ok := c.users[username]; !ok {
		ErrUserNotFound.ErrInfo = username
		return nil, ErrUserNotFound
	}
	token, ok := tokens[username]
	if ok {
		return token, nil
	} 
	ErrJWTNotYetFetched.ErrInfo = username
	return nil, ErrJWTNotYetFetched
}

func (c *JWTClient) FetchAllTokens(username string) error {

	for user, userCred := c.users {
		token, err := fetchToken(userCred)
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

	jwtUrl, err := resolvePath(c.BaseURL, "token")
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
		c.httpClientP = httpClient()
	}
	rsp, err := c.httpClientP.Do(req)
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode != http.StatusOK {
		ErrFailedToFetchJWTToken.ErrCode = rsp.StatusCode
		return nil, ErrFailedToFetchJWTToken
	}
	jwtToken, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	return jwtToken, nil
}
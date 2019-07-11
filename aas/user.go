package aas

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	types "intel/isecl/lib/common/types/aas"
)

type UserClientErr struct {
	ErrMessage string
	ErrCode    int
}

func (ucErr *UserClientErr) Error() string {
	return fmt.Sprintf("%s: http status %d", ucErr.ErrMessage, ucErr.ErrCode)
}

var (
	ErrFailedToGetJWTCert  = &UserClientErr{"failed to retrieve JWT signing certificate", 0}
	ErrFailedToGetJWTToken = &UserClientErr{"failed to retrieve JWT token", 0}
)

type UserClient struct {
	BaseURL  string
	Username string
	Password string
	JWTToken []byte

	HTTPClient *http.Client

	userCred *types.UserCred
}

func (c *UserClient) GetJwtSigningCert() ([]byte, error) {

	jwtCertUrl, err := resolvePath(c.BaseURL, "noauth/jwtCert")
	if err != nil {
		return nil, err
	}
	req, _ := http.NewRequest(http.MethodGet, jwtCertUrl, nil)
	req.Header.Set("Accept", "application/x-pem-file")

	if c.HTTPClient == nil {
		c.HTTPClient = &http.Client{}
	}
	rsp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode != http.StatusOK {
		ErrFailedToGetJWTCert.ErrCode = rsp.StatusCode
		return nil, ErrFailedToGetJWTCert
	}
	return ioutil.ReadAll(rsp.Body)
}

func (c *UserClient) GetJwtToken() ([]byte, error) {

	jwtUrl, err := resolvePath(c.BaseURL, "token")
	if err != nil {
		return nil, err
	}
	if c.userCred == nil {
		c.userCred = &types.UserCred{
			UserName: c.Username,
			Password: c.Password,
		}
	}
	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(c.userCred)
	if err != nil {
		return nil, err
	}
	req, _ := http.NewRequest("POST", jwtUrl, buf)
	req.Header.Set("Accept", "application/jwt")

	if c.HTTPClient == nil {
		c.HTTPClient = &http.Client{}
	}
	rsp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode != http.StatusOK {
		ErrFailedToGetJWTToken.ErrCode = rsp.StatusCode
		return nil, ErrFailedToGetJWTToken
	}
	return ioutil.ReadAll(rsp.Body)
}

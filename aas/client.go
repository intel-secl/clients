package aas

import (
	"bytes"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"fmt"
	"intel/isecl/lib/clients"
	types "intel/isecl/lib/common/types/aas"
)

type Client struct {
	BaseURL    string
	JWTToken   []byte
	HTTPClient *http.Client
}

var (
	ErrHTTPCreateUser = &clients.HTTPClientErr{
		ErrMessage: "Failed to create user",
		ErrCode:    0,
	}
	ErrHTTPCreateRole = &clients.HTTPClientErr{
		ErrMessage: "Failed to create role",
		ErrCode:    0,
	}
	ErrHTTPAddRoleToUser = &clients.HTTPClientErr{
		ErrMessage: "Failed to add role to user",
		ErrCode:    0,
	}
)

func (c *Client) prepReqHeader(req *http.Request) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+string(c.JWTToken))
}

func (c *Client) CreateUser(u types.UserCreate) (*types.UserCreateResponse, error) {

	userURL := clients.ResolvePath(c.BaseURL, "users")

	payload, err := json.Marshal(&u)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, userURL, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	c.prepReqHeader(req)

	if c.HTTPClient == nil {
		c.HTTPClient = clients.HTTPClientTLSNoVerify()
	}
	rsp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode != http.StatusOK {
		ErrHTTPCreateUser.ErrCode = rsp.StatusCode
		return nil, ErrHTTPCreateUser
	}
	var userCreateResponse types.UserCreateResponse
	err = json.NewDecoder(rsp.Body).Decode(&userCreateResponse)
	if err != nil {
		return nil, err
	}
	return &userCreateResponse, nil
}

func (c *Client) CreateRole(r types.RoleCreate) (*types.RoleCreateResponse, error) {

	roleURL := clients.ResolvePath(c.BaseURL, "roles")

	payload, err := json.Marshal(&r)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, roleURL, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	c.prepReqHeader(req)

	if c.HTTPClient == nil {
		c.HTTPClient = clients.HTTPClientTLSNoVerify()
	}
	rsp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	msg, _ := ioutil.ReadAll(rsp.Body)
	if rsp.StatusCode != http.StatusCreated {
		ErrHTTPCreateRole.ErrCode = rsp.StatusCode
		return nil, fmt.Errorf("Failed to create role: HTTP Code: %d: Response Message: %s", ErrHTTPCreateRole.ErrCode, msg)
	}
	var roleCreateResponse types.RoleCreateResponse
	err = json.NewDecoder(rsp.Body).Decode(&roleCreateResponse)
	if err != nil {
		return nil, err
	}
	return &roleCreateResponse, nil
}

func (c *Client) AddRoleToUser(userID string, r types.UserRoleCreate) error {

	userRoleURL := clients.ResolvePath(c.BaseURL, "users"+userID+"roles")

	payload, err := json.Marshal(&r)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, userRoleURL, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	c.prepReqHeader(req)

	if c.HTTPClient == nil {
		c.HTTPClient = clients.HTTPClientTLSNoVerify()
	}
	rsp, err := c.HTTPClient.Do(req)
	if rsp.StatusCode != http.StatusOK {
		ErrHTTPAddRoleToUser.ErrCode = rsp.StatusCode
		return ErrHTTPAddRoleToUser
	}
	return nil
}

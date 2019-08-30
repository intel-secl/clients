/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package aas

import (
	"bytes"
	"encoding/json"
	"intel/isecl/lib/clients"
	types "intel/isecl/lib/common/types/aas"
	"io/ioutil"
	"net/http"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	BaseURL    string
	JWTToken   []byte
	HTTPClient *http.Client
}

var (
	ErrHTTPCreateUser = &clients.HTTPClientErr{
		ErrMessage: "Failed to create user",
	}
	ErrHTTPCreateRole = &clients.HTTPClientErr{
		ErrMessage: "Failed to create role",
	}
	ErrHTTPAddRoleToUser = &clients.HTTPClientErr{
		ErrMessage: "Failed to add role to user",
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
		ErrHTTPCreateUser.RetCode = rsp.StatusCode
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
		ErrHTTPCreateRole.RetCode = rsp.StatusCode
		ErrHTTPCreateRole.RetMessage = string(msg)
		log.Errorf("Role not created. http errorcode : %d, message: %s", ErrHTTPCreateRole.RetCode, ErrHTTPCreateRole.RetMessage )
		return nil, ErrHTTPCreateRole
	}
	var roleCreateResponse types.RoleCreateResponse
	err = json.Unmarshal(msg, &roleCreateResponse)
	if err != nil {
		log.WithError(err).Error("CreateRole could not decode response")
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
		ErrHTTPAddRoleToUser.RetCode = rsp.StatusCode
		return ErrHTTPAddRoleToUser
	}
	return nil
}

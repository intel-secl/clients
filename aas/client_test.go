/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package aas

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"

	types "intel/isecl/lib/common/v5/types/aas"
)

// Run with command: go test -count=1 -v <filenames>
// insert the url to a working aas instance to "aasURL" variable
func TestAASClient(t *testing.T) {

	var token []byte
	var err error

	aasMockSrv, port := aasMockServer(t)
	defer aasMockSrv.Close()

	aasURL := "http://localhost" + port + "/aas"

	// get token of aas admin
	jwt := NewJWTClient(aasURL)
	jwt.HTTPClient = http.DefaultClient
	jwt.AddUser("admin", "password")
	err = jwt.FetchAllTokens()

	assert.NoError(t, err, "All user tokens should be fetched")

	token, err = jwt.GetUserToken("admin")
	fmt.Println("token: ", string(token))

	assert.NoError(t, err, "admin token should be fetched")

	aasClient := Client{
		BaseURL:  aasURL,
		JWTToken: token,
	}
	aasClient.HTTPClient = http.DefaultClient
	role := types.RoleCreate{
		RoleInfo: types.RoleInfo{
			Service: "test_service",
			Name:    "test_name",
		},
	}
	_, err = aasClient.CreateRole(role)
	assert.NoError(t, err, "role should be created")
}

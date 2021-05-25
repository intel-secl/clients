/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package aas

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	types "intel/isecl/lib/common/v4/types/aas"
	"net"
	"net/http"
	"testing"
)

const (
	aasToken       = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJyb2xlcyI6W3sic2VydmljZSI6IkFBUyIsIm5hbWUiOiJSb2xlTWFuYWdlciJ9LHsic2VydmljZSI6IkFBUyIsIm5hbWUiOiJVc2VyTWFuYWdlciJ9LHsic2VydmljZSI6IkFBUyIsIm5hbWUiOiJVc2VyUm9sZU1hbmFnZXIifSx7InNlcnZpY2UiOiJUQSIsIm5hbWUiOiJBZG1pbmlzdHJhdG9yIn0seyJzZXJ2aWNlIjoiVlMiLCJuYW1lIjoiQWRtaW5pc3RyYXRvciJ9LHsic2VydmljZSI6IktNUyIsIm5hbWUiOiJLZXlDUlVEIn0seyJzZXJ2aWNlIjoiQUgiLCJuYW1lIjoiQWRtaW5pc3RyYXRvciJ9LHsic2VydmljZSI6IldMUyIsIm5hbWUiOiJBZG1pbmlzdHJhdG9yIn1dLCJwZXJtaXNzaW9ucyI6W3sic2VydmljZSI6IkFIIiwicnVsZXMiOlsiKjoqOioiXX0seyJzZXJ2aWNlIjoiS01TIiwicnVsZXMiOlsiKjoqOioiXX0seyJzZXJ2aWNlIjoiVEEiLCJydWxlcyI6WyIqOio6KiJdfSx7InNlcnZpY2UiOiJWUyIsInJ1bGVzIjpbIio6KjoqIl19LHsic2VydmljZSI6IldMUyIsInJ1bGVzIjpbIio6KjoqIl19XSwiZXhwIjoxNjA2Mjg1MDA0LCJpYXQiOjE1OTQ0NzQwMDEsImlzcyI6IkFBUyBKV1QgSXNzdWVyIiwic3ViIjoiYWRtaW4iLCJqdGkiOiJmOTBlZGU4YS00MzU5LTQyZTktOTU0ZS0wNDA5MmI4YmE3YjQifQ.g_UFb8xHLIL8nFNQ3XGE2ne4Eic0MFTdcZ_dQjoCDMQ"
	jwtCertsB64Enc = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUVBekNDQW11Z0F3SUJBZ0lCQmpBTkJna3Foa2lHOXcwQkFRd0ZBREJRTVFzd0NRWURWUVFHRXdKVlV6RUwKTUFrR0ExVUVDQk1DVTBZeEN6QUpCZ05WQkFjVEFsTkRNUTR3REFZRFZRUUtFd1ZKVGxSRlRERVhNQlVHQTFVRQpBeE1PUTAxVElGTnBaMjVwYm1jZ1EwRXdIaGNOTWpBeE1URXdNRE14TURNNVdoY05NakV4TVRFd01ETXhNRE01CldqQW1NU1F3SWdZRFZRUURFeHRCUVZNZ1NsZFVJRk5wWjI1cGJtY2dRMlZ5ZEdsbWFXTmhkR1V3Z2dHaU1BMEcKQ1NxR1NJYjNEUUVCQVFVQUE0SUJqd0F3Z2dHS0FvSUJnUUN2YjZOTXBQcEozMWhjY015Q0w3SUNQQWw1eDBCeAp2S0xDZWFJRy9WVnhqWlN1SlJ6MVJmdFJNR1RsS3hRY2p6TkNSNUdZSVA5T2ZHY0RUUzNnR1RvY1YwdzhhSVNnCnVhdG1ReHhPSEQ5blpUNjhsT0ZNT3JUM3dzMUpUY01Xd2kxeXBQb1NEdHVjbXBQSmJiVU1VVmMyN3ZUUFcwZk4KYjVBcnNMODBUVnhwSG9JZ0Y2ZUtFV29TZStFcEM2N2ZBN0lpaGl0NzY4U2pIMUNEQkxlOU1HNDlKSFlac2JNZgpQd2RtSTF4UjBmeS9raGRhcXVDZFBZenJWT0xjcnpGeFk2cnUrMVQ1N3JJYzFXcVhvL2RCWWNoMkF2TEQvbzdaCm1wcU9YWWRQWll4QmFMK2NMWHk0SGNKaG5XRU1BWDVjVWVqcWYrQTlzMWROcE9VZEpwRnhUdUpraGRuRERKVnUKOS9odVZoaWVxVXplOFY2aytnN3grUjhPWVNzMDZaSnovMEg3V2I0VXgwMDErRFI2SjZ4cjh5RFd0T3ZFcm8vdwpDQWJZdUJxN1VtZGw1NXdIdFE0dkhlaDhxY0ZYQ1pJbVpLTjYyTkZ6a1JnOExCTTAzQUxQWFZBR2ZwODlnNEx6CnRscXkvbXNVbmxrZ2JBVks0cjJCMFpIamhzN2doTVl4SVE4Q0F3RUFBYU1TTUJBd0RnWURWUjBQQVFIL0JBUUQKQWdiQU1BMEdDU3FHU0liM0RRRUJEQVVBQTRJQmdRQU4xUGlWeXVRNzU4QVAvVlR6S1BnS2ZlNCtZNHduQ2dycwpObUhKeHIrWGhFYjh3TzdMa3FXTWRMOTN0NWZ5REcrRXJXcEV6d2FUMDUxQThLbGJ6MjVNd0NjTncvb2RTRDNNCkQxdGY5L2ZMVVN0YVR1RjRpUEZpdnRpdWxZMmZKWC8vM0thbStsRHhYWndpZDlKemNPbDV6bGhhL0lZVFlLRlMKdzUwdVN0TlN2VXFPVjB6Q1BaQ0NPZ1BkeWM5SGs0OHI0V2FkZE96QlZ0a1NDMmd5SFA1cnorcE5GYUgwckJhMgorRk1pVjkvQkJkckxReG83dWxiVXRKRDJFYkdZdHlMQ1NQV3pGK244T0xMem16SmhNMjBNVWd6c3pvdnJSaGRkCi85bENjcU5mNzhjUFBHa1gzcWdUdVlyb2tYWlZmVHJ6NVpTVFRzOVpnd2pMbmk1YUhmSytUNmI3S09URXgzVm4KUG03aG9oeEduREdEK0tBZlZuVnZ4TDlWMEJLcTl6VVQ1aEYwa0FjS3lMdWk0YVVCcWxGeStQMm9RWVlhc293VgpmUWx6dE5kdjhqbzBFQ2F3RDVQcWQrZTZBd2dsakorWm5PZW0rakVhMGE2b0xBUDh4bFA3Ukp1UTFlUzZCdnhjCnVpTnNkalNqU0VJYUlockZEUXA0WDI5bXVoa0hlTG89Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0KLS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUVOVENDQXAyZ0F3SUJBZ0lCQXpBTkJna3Foa2lHOXcwQkFRd0ZBREJITVFzd0NRWURWUVFHRXdKVlV6RUwKTUFrR0ExVUVDQk1DVTBZeEN6QUpCZ05WQkFjVEFsTkRNUTR3REFZRFZRUUtFd1ZKVGxSRlRERU9NQXdHQTFVRQpBeE1GUTAxVFEwRXdIaGNOTWpBeE1URXdNRE14TURNeFdoY05NalV4TVRFd01ETXhNRE14V2pCUU1Rc3dDUVlEClZRUUdFd0pWVXpFTE1Ba0dBMVVFQ0JNQ1UwWXhDekFKQmdOVkJBY1RBbE5ETVE0d0RBWURWUVFLRXdWSlRsUkYKVERFWE1CVUdBMVVFQXhNT1EwMVRJRk5wWjI1cGJtY2dRMEV3Z2dHaU1BMEdDU3FHU0liM0RRRUJBUVVBQTRJQgpqd0F3Z2dHS0FvSUJnUURYaFVySFVGbzE3QTdkOWpDUVVVTmprS0hVYlluK0EwaUlCTERNY0ZCNE8xcFYrMVh5CnRXMzE3Q2JuN05qRkY5ajA0TURVeWFLKzZKR1pxNEZ4VU5HeFF1QXlZWk9lVTFlV2NXczM5Z0JWemh3Z0s0N00KOVE4dnFHOGNITW10OGZiYUdjc0JCelM3dDRFbVNJNWtWVzJUNWVtSkE5MlhsVVEveWZhSmlGZDNLU3l5MDB6KwpKKzByVE4xT1ZTYWpsa1BKTm5hWCtQTmVzNEtLVGVsVEtZODRKMlBIaHhjNVRuRmRidTFBYlZ4Zkhlb1NWQ0hGCmw5R3NMWkxpcEJnMDZoWHBtSGxhd3Bzb2Z6QWFBNE9FN0xGZVlLSnMwZUF6dzVWcE4zV2hGemNrdFZxcDg1SnAKbXVTWGZYTTJRQUw5d0U1Ri8wSXpLMU9LR3pRRS9MaFFvVXNxYUd6bGZYc01IRjFPWjZJVGZKaG9YK1JMWEtYNApmWVo3Wi9xeE5SK2plQmxueW5FS0J1TFlVdThkRHN4WHAvTU8zdFFMU3BGL3M5dUc5WDZxeUxwVGQ3ckYxWlMrCm9mR1NNME9aTEJySHhEVkErbnhVVE40clg1My96Ymg2MVNWT2dIalJMQ3cwNXFYN3RncVhiQVlzQklFNDBqeTIKaG9lYUJiNGJkOWdzSHBzQ0F3RUFBYU1qTUNFd0RnWURWUjBQQVFIL0JBUURBZ0VHTUE4R0ExVWRFd0VCL3dRRgpNQU1CQWY4d0RRWUpLb1pJaHZjTkFRRU1CUUFEZ2dHQkFCKzdUVUlhTXBtUmp6R2Y4Ni8weDRBYUNHZHFmTmJkCkZ1N3RLZ0hIVlZYNGV5ZlNhNitBcEE5bnl3R2s5RVY5Q0FLenVRL0Q5R1g2ZHdSYjd2SVhNai80WklscXdLbUoKUVBjTm9CYm10Y2lnbXB3SmM2c016elBZRmFMRHAxNEVIMklvRldjNEVYRzNqT0E1ZkNXZktad3owT0JpcUFSegprd2hOTmF2dWRxR09pUldWNy9yVGw0aEh5VG1hdHUzaStOdlh2L3JSUVljeEwzZUkwNXUvWlpPTnZWdDhlajd0CnRaM0RDZ2UwUUhQS1ppK0Y3dkMwUWRjQVZjUEdEUHM3ZUFieUxqQlgrUGRPc3g5N3owQzBqL093VGdRbHQxNEkKMG00S0w4bXViNEVwcW9EN0dMaDFXN0xvaFp5RXNLeW1LR0I1T3JLWHJPVnBPVXFrMGJaUlBQOGx2WS8zWkxaUQpkYVdRUENRdjUybHVMRjVxR0tCbFFqZWw4NlB6K2tKczhvd2I3eUwxTmxaRktrTWJlK3JhTXlESUU5UVZzWUw3CkMrVHVUTGFuY2xtbjdJN2I1Nm1KcklXNWlXL3BDeDZFK0c2NGs5Z0dHdUxCQjk3RzdwNEM0TDg2QUU3SXltUSsKdzZKb3V6SGhiQUdCQVlDdU1lZW0wOXFQM0VpV01BQjk2UT09Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0="
)

func tokenMockGoodResponse(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/jwt")
	_, _ = w.Write([]byte(aasToken))
}

func roleCreateMockGoodResponse(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	roleCResp := types.RoleCreateResponse{
		Service: "test_service",
		Name:    "test_id",
		ID:      uuid.New().String(),
	}
	resp, _ := json.Marshal(roleCResp)
	_, _ = w.Write(resp)
}
func aasMockServer(t *testing.T) (*http.Server, string) {
	handler := mux.NewRouter()
	handler.HandleFunc("/aas/noauth/jwt-certificates", mockJwtSigningCertResponse).Methods("GET")
	handler.HandleFunc("/aas/token", tokenMockGoodResponse).Methods("POST")
	handler.HandleFunc("/aas/roles", roleCreateMockGoodResponse).Methods("POST")

	return mockServerLauncher(t, handler)
}

// mockServerLauncher serves the mock server on a dynamically allocated port
func mockServerLauncher(t *testing.T, r http.Handler) (*http.Server, string) {

	//Listener Implementations
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Log("mockServerLauncher() : Unable to initiate Listener", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	err = listener.Close()
	if err != nil {
		t.Log("mockServerLauncher() : Unable to close Listener", err)
	}
	portString := fmt.Sprintf(":%d", port)

	h := &http.Server{
		Addr:    portString,
		Handler: r,
	}
	go h.ListenAndServe()

	return h, portString
}

func mockJwtSigningCertResponse(w http.ResponseWriter, r *http.Request) {
	resp, _ := base64.StdEncoding.DecodeString(jwtCertsB64Enc)

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/x-pem-file")
	_, _ = w.Write(resp)
}

// Run with command: go test -count=1 -v <filenames>
func TestJWT(t *testing.T) {

	aasMockSrv, port := aasMockServer(t)
	defer aasMockSrv.Close()

	aasURL := "http://localhost" + port + "/aas"

	jwt := NewJWTClient(aasURL)
	jwt.HTTPClient = http.DefaultClient

	cert, certErr := jwt.GetJWTSigningCert()
	fmt.Println(cert, certErr)

	var token []byte
	var err error
	jwt.AddUser("user1", "password")
	jwt.AddUser("user2", "password")

	fmt.Println("fetch 1")
	token, err = jwt.FetchTokenForUser("user1")
	assert.NoError(t, err, "FetchTokenForUser should be successful")

	fmt.Println("fetch all")
	err = jwt.FetchAllTokens()
	assert.NoError(t, err, "FetchAllTokens  should be successful")

	token, err = jwt.GetUserToken("user1")
	t.Log("token 1: ", string(token))
	assert.NoError(t, err, "user1 token should be present")

	token, err = jwt.GetUserToken("user2")
	t.Log("token 2: ", string(token))
	assert.NoError(t, err, "user2 token should be present")

	_, err = jwt.GetUserToken("user3")
	assert.Error(t, err, "user3 token should NOT be present")

	jwt404 := NewJWTClient("https://url.to.aas.instance:port/abc")
	jwt.HTTPClient = http.DefaultClient
	_, err = jwt404.GetJWTSigningCert()
	assert.Error(t, err, "Signing certificate fetch should return 404")

	jwt404.AddUser("user1", "password")
	err = jwt404.FetchAllTokens()
	assert.Error(t, err, "User token fetch should fail")
}

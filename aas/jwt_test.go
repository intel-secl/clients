package aas

import (
	"fmt"
	"testing"
)

func TestJWT(t *testing.T) {

	jwt := JWTClient{
		BaseURL:        "https://10.105.168.81:8443/aas",
		NoTLSKeyVerify: true,
	}
	cert, certErr := jwt.GetJWTSigningCert()
	fmt.Println(cert, certErr)

	var token []byte
	var err error
	jwt.AddUser("user1", "password")
	jwt.AddUser("user2", "password")
	jwt.AddUser("user3", "password")

	fmt.Println("fetch 1")
	jwt.FetchTokenForUser("user1")
	token, err = jwt.GetUserToken("user1")
	fmt.Println("token 1: ", token)
	fmt.Println("err1 1: ", err.Error())
	token, err = jwt.GetUserToken("user2")
	fmt.Println("token 2")
	fmt.Println("err1 2: ", err.Error())

	fmt.Println("fetch all")
	jwt.FetchAllTokens()

	token, err = jwt.GetUserToken("user1")
	fmt.Println("token 1: ", token)
	fmt.Println("err1 1: ", err.Error())
	token, err = jwt.GetUserToken("user2")
	fmt.Println("token 2")
	fmt.Println("err1 2: ", err.Error())
	token, err = jwt.GetUserToken("user3")
	fmt.Println("token 3")
	fmt.Println("err1 3: ", err.Error())
	token, err = jwt.GetUserToken("user4")
	fmt.Println("err1 4: ", err.Error())
}

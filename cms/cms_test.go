package cms

import (
	"fmt"
	"io/ioutil"
	"testing"
	// "github.com/stretchr/testify/assert"
)

func TestCMS(t *testing.T) {

	cms := Client{
		BaseURL: "",
	}
	jwtToken, err := ioutil.ReadFile("/var/jwtToken")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	cms.JWTToken = jwtToken

}

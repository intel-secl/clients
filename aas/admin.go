package aas

import (
	"fmt"
	types "intel/isecl/lib/common/types/aas"
	"net/http"
)

type AdminClientErr struct {
	ErrMessage string
	ErrCode    int
}

func (ucErr *AdminClientErr) Error() string {
	return fmt.Sprintf("%s: http error %d", ucErr.ErrMessage, ucErr.ErrCode)
}

var (
	ErrAdmin = &AdminClientErr{ErrMessage: "failed to retrieve JWT signing certificate"}
)

type AdminClient struct {
	BaseURL  string
	Username string
	Password string
	JWTToken []byte

	HTTPClient *http.Client

	userCred types.UserCred
}

// func (c *Client) CreateUser(u sharedTypes.UserCreate) (*sharedTypes.UserCreateResponse, error) {
// 	userURL, err := c.resolvePath("users")
// 	if err != nil {
// 		return nil, err
// 	}

// 	payload, err := json.Marshal(&u)
// 	if err != nil {
// 		return nil, err
// 	}

// 	req, err := http.NewRequest(http.MethodPost, userURL, bytes.NewBuffer(payload))
// 	if err != nil {
// 		return nil, errors.New("Error creating request")
// 	}
// 	req.Header.Set("Accept", "application/json")
// 	req.Header.Set("Content-Type", "application/json")

// 	rsp, err := c.dispatchRequest(req)
// 	if err != nil {
// 		return nil, fmt.Errorf("Failed to create dispatch request: %v", err)
// 	}

// 	if rsp.StatusCode != http.StatusCreated {
// 		return nil, fmt.Errorf("Failed to create user: HTTP Code: %d", rsp.StatusCode)
// 	}
// 	var userCreateResponse sharedTypes.UserCreateResponse
// 	err = json.NewDecoder(rsp.Body).Decode(&userCreateResponse)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &userCreateResponse, nil
// }

// // CreateRole is used to create role in aas
// func (c *Client) CreateRole(r sharedTypes.RoleCreate) (*sharedTypes.RoleCreateResponse, error) {

// 	roleURL, _ := c.resolvePath("roles")

// 	payload, err := json.Marshal(&r)
// 	if err != nil {
// 		return nil, err
// 	}

// 	req, err := http.NewRequest(http.MethodPost, roleURL, bytes.NewBuffer(payload))
// 	if err != nil {
// 		return nil, errors.New("Error creating request")
// 	}
// 	req.Header.Set("Accept", "application/json")
// 	req.Header.Set("Content-Type", "application/json")

// 	rsp, err := c.dispatchRequest(req)
// 	if err != nil {
// 		return nil, fmt.Errorf("Failed to create dispatch request: %v", err)
// 	}
// 	if rsp.StatusCode != http.StatusCreated {
// 		return nil, fmt.Errorf("Failed to create role: HTTP Code: %d", rsp.StatusCode)
// 	}
// 	var roleCreateResponse sharedTypes.RoleCreateResponse
// 	err = json.NewDecoder(rsp.Body).Decode(&roleCreateResponse)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &roleCreateResponse, nil
// }

// // CreateRole is used to create role in aas
// func (c *Client) AddRoleToUser(userID string, r sharedTypes.UserRoleCreate) error {
// 	reqURL, err := c.resolvePath("users" + userID + "roles")
// 	if err != nil {
// 		return err
// 	}

// 	payload, err := json.Marshal(&r)
// 	if err != nil {
// 		return err
// 	}

// 	req, err := http.NewRequest(http.MethodPost, reqURL, bytes.NewBuffer(payload))
// 	if err != nil {
// 		return errors.New("Error creating request")
// 	}
// 	req.Header.Set("Accept", "application/json")
// 	req.Header.Set("Content-Type", "application/json")

// 	rsp, err := c.dispatchRequest(req)
// 	if rsp.StatusCode != http.StatusOK {
// 		return fmt.Errorf("Failed to add role user: HTTP Code: %d", rsp.StatusCode)
// 	}

// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

package aas

import (
	"fmt"
	types "intel/isecl/lib/common/types/aas"
	"net/http"
)

type Client struct {
	BaseURL  string
	JWTToken []byte
	httpClientP *http.Client
}

// func (c *AdminClient) prepRequest(req *http.Request) (*http.Response, error) {

// 	if c.userCred == nil {
// 		c.userCred = &types.UserCred{
// 			UserName: c.Username,
// 			Password: c.Password,
// 		}
// 	}
// }

// func (c *AdminClient) CreateUser(u types.UserCreate) (*types.UserCreateResponse, error) {
// 	userURL, err := resolvePath(c.BaseURL, "users")
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
// 	var userCreateResponse types.UserCreateResponse
// 	err = json.NewDecoder(rsp.Body).Decode(&userCreateResponse)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &userCreateResponse, nil
// }

// // CreateRole is used to create role in aas
// func (c *AdminClient) CreateRole(r types.RoleCreate) (*types.RoleCreateResponse, error) {

// 	roleURL, _ := resolvePath(c.BaseURL, "roles")

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
// 	var roleCreateResponse types.RoleCreateResponse
// 	err = json.NewDecoder(rsp.Body).Decode(&roleCreateResponse)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &roleCreateResponse, nil
// }

// // CreateRole is used to create role in aas
// func (c *AdminClient) AddRoleToUser(userID string, r types.UserRoleCreate) error {
// 	reqURL, err := resolvePath(c.BaseURL, "users"+userID+"roles")
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

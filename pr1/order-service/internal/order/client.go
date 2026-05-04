package order

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type UserServiceClient struct {
	baseURL string
	client  *http.Client
}

func NewUserServiceClient(baseURL string) *UserServiceClient {
	return &UserServiceClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 3 * time.Second},
	}
}

func (c *UserServiceClient) GetUserByID(id int64) (UserDTO, error) {
	url := fmt.Sprintf("%s/users/%d", c.baseURL, id)

	resp, err := c.client.Get(url)
	if err != nil {
		return UserDTO{}, fmt.Errorf("request to user-service failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return UserDTO{}, fmt.Errorf("user not found in user-service")
	}
	if resp.StatusCode != http.StatusOK {
		return UserDTO{}, fmt.Errorf("user-service returned status %d", resp.StatusCode)
	}

	var user UserDTO
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return UserDTO{}, fmt.Errorf("decode user response failed: %w", err)
	}
	return user, nil
}

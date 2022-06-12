package usersvc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"path"
	"strings"
)

type HTTPClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewHTTPClient(httpClient *http.Client, baseURL string) (*HTTPClient, error) {
	return &HTTPClient{
		httpClient: httpClient,
		baseURL:    strings.TrimSuffix(baseURL, "/"),
	}, nil
}

func (c *HTTPClient) GetUser(ctx context.Context, name string) (user *User, err error) {
	url := c.baseURL + path.Join("/users", name)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, jsonToError(resp.Body)
	}

	user = new(User)
	if err := json.NewDecoder(resp.Body).Decode(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (c *HTTPClient) ListUsers(ctx context.Context) (users []*User, err error) {
	url := c.baseURL + "/users"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, jsonToError(resp.Body)
	}

	var out = struct {
		Users []*User `json:"users"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out.Users, nil
}

func (c *HTTPClient) CreateUser(ctx context.Context, user *User) (err error) {
	url := c.baseURL + "/users"
	body, err := json.Marshal(user)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return jsonToError(resp.Body)
	}

	return nil
}

func (c *HTTPClient) UpdateUser(ctx context.Context, name string, user *User) (err error) {
	url := c.baseURL + path.Join("/users", name)
	body, err := json.Marshal(user)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return jsonToError(resp.Body)
	}

	return nil
}

func (c *HTTPClient) DeleteUser(ctx context.Context, name string) (err error) {
	url := c.baseURL + path.Join("/users", name)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return jsonToError(resp.Body)
	}

	return nil
}

func jsonToError(r io.Reader) error {
	var resp map[string]string
	if err := json.NewDecoder(r).Decode(&resp); err != nil {
		return err
	}

	errorStr := resp["error"]
	if errorStr == "" {
		return nil
	}
	return errors.New(errorStr)
}

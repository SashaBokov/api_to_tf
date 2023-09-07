package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"net/http"
	"net/url"
)

const (
	// endpoints
	// runtimeGroupEndpoint is the endpoint for operations with a runtime group.
	runtimeGroupEndpoint = "/create-runtime-group"

	// methods
	// createRuntimeGroupMethod is the HTTP method for creating a runtime group.
	createRuntimeGroupMethod = http.MethodPost
)

// Client is the representation of http client for the GroupAPI.
type Client struct {
	BaseUrl string
	token   string
}

// New is a constructor for Client.
func New(baseULR, token string) (*Client, error) {
	client := &Client{}

	// baseULR validation.
	_, err := url.Parse(baseULR)
	if err != nil {
		return nil, client.wrap("error parsing base URL", err)
	}

	// token validation.
	if err := validateBearerToken(token); err != nil {
		return nil, client.wrap("error validating bearer token", err)
	}

	client.BaseUrl = baseULR
	client.token = token

	return client, nil
}

// CreateRuntimeGroupRequest represents the request body for creating a runtime group.
type CreateRuntimeGroupRequest struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	ClusterType string            `json:"cluster_type"`
	Labels      map[string]string `json:"labels"`
}

// CreateRuntimeGroupResponse represents the response from creating a runtime group.
type CreateRuntimeGroupResponse struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Labels      map[string]string `json:"labels"`
	Config      struct {
		ControlPlaneEndpoint string `json:"control_plane_endpoint"`
		TelemetryEndpoint    string `json:"telemetry_endpoint"`
	} `json:"config"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// CreateRuntimeGroup sends a POST request to create a runtime group.
func (c *Client) CreateRuntimeGroup(requestBody CreateRuntimeGroupRequest) (*CreateRuntimeGroupResponse, error) {
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, c.wrap(" serializing request body", err)
	}

	endpoint, err := url.JoinPath(c.BaseUrl, runtimeGroupEndpoint)
	if err != nil {
		return nil, c.wrap(" joining base URL and endpoint", err)
	}

	req, err := http.NewRequest(createRuntimeGroupMethod, endpoint, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return nil, c.wrap("creating HTTP request", err)
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, c.wrap("making HTTP request", err)
	}
	defer resp.Body.Close()

	// Check the HTTP response status code.
	if err := c.codeToErr(resp.StatusCode); err != nil {
		return nil, c.wrap("checking status code", err)
	}

	var createResponse CreateRuntimeGroupResponse
	if err := json.NewDecoder(resp.Body).Decode(&createResponse); err != nil {
		return nil, c.wrap("decoding response JSON", err)
	}

	return &createResponse, nil
}

// do is a wrapper for http.Client.Do
func (c *Client) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	// Perform the HTTP request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, c.wrap("making HTTP request", err)
	}

	return resp, nil
}

// wrap the client function for wrapping the error.
func (c *Client) wrap(msg string, err error) error {
	return fmt.Errorf("|client error: %s -> %w", msg, err)
}

func (c *Client) codeToErr(code int) error {
	if code != http.StatusCreated && code != http.StatusOK {
		// todo: Handle error responses (e.g., 400, 401, 403, 409, 500, 503)
		return fmt.Errorf("HTTP request failed with status code %d", code)
	}
	return nil
}

// ValidateBearerToken validates the bearer token.
func validateBearerToken(tokenString string) error {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return err
	}

	// Check if the token is valid.
	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

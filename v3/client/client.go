package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"go.gh.ink/openapi/sdk/20260512/v3"
	"go.gh.ink/openapi/sdk/20260512/v3/errors"
)

// Client provides basic struct for client object
type Client struct {
	endpoint           string
	secretID           string
	secretKey          string
	enableToken        bool
	token              string
	timeout            int
	maxRetries         int
	retryDelay         int
	exponentialBackoff bool
	marshal            func(any) ([]byte, error)
	unmarshal          func([]byte, any) error
	Logger             Logger
}

// Option provides a basic option type
type Option func(*Client)

// WithLogger sets default logger to custom
func WithLogger(logger Logger) Option {
	return func(c *Client) {
		c.Logger = logger
	}
}

// WithEndpoint sets default endpoint
func WithEndpoint(endpoint string) Option {
	return func(c *Client) {
		c.endpoint = endpoint
	}
}

// WithMarshal sets default marshal lib
func WithMarshal(marshal func(any) ([]byte, error)) Option {
	return func(c *Client) {
		c.marshal = marshal
	}
}

// WithUnmarshal sets default unmarshal lib
func WithUnmarshal(unmarshal func([]byte, any) error) Option {
	return func(c *Client) {
		c.unmarshal = unmarshal
	}
}

// WithTimeout sets timeout for request
func WithTimeout(timeout int) Option {
	return func(c *Client) {
		c.timeout = timeout
	}
}

// WithMaxRetries sets max retries for request
// If maxRetries <= 0, defaults to DefaultMaxRetries (5 attempts)
func WithMaxRetries(maxRetries int) Option {
	if maxRetries <= 0 {
		maxRetries = DefaultMaxRetries
	}
	return func(c *Client) {
		c.maxRetries = maxRetries
	}
}

// WithRetryDelay sets initial retry delay for request in seconds
// The actual delay grows exponentially if exponential backoff is enabled,
// but is capped at MaxRetryDelaySeconds (60 seconds)
// If retryDelay <= 0, defaults to DefaultRetryDelaySeconds (1 second)
func WithRetryDelay(retryDelay int) Option {
	if retryDelay <= 0 {
		retryDelay = DefaultRetryDelaySeconds
	}
	return func(c *Client) {
		c.retryDelay = retryDelay
	}
}

// WithExponentialBackoff enables/disables exponential backoff for request retries
// When enabled (true), retry delays grow exponentially: delay = delay * 2
// Delays are capped at MaxRetryDelaySeconds (60 seconds)
// By default, exponential backoff is enabled (DefaultExponentialBackoff = true)
func WithExponentialBackoff(exponentialBackoff bool) Option {
	return func(c *Client) {
		c.exponentialBackoff = exponentialBackoff
	}
}

// EnableToken enables token as authorisation
func EnableToken(enableToken bool) Option {
	return func(c *Client) {
		c.enableToken = enableToken
	}
}

// GetEndpoint returns endpoint
func (c *Client) GetEndpoint() string {
	return c.endpoint
}

// applyToken applies a new token
func applyToken(c *Client) error {
	// Send request
	result := c.Send(
		strings.Join([]string{c.endpoint, "/openapi/token"}, ""),
		http.MethodGet,
		nil,
	).WithKey()
	if result.Err != nil {
		c.Logger.Error(nil, fmt.Sprintf(
			"failed to get token, sender error: %s", result.Err.Error(),
		))
		return errors.ErrTokenAcquisitionFailed.
			WithApiCode(result.Code).
			WithApiMessage(result.Msg).
			WithRequestID(result.RequestID).
			WithResponse(result)
	}

	// Check status code
	if !result.OK() {
		c.Logger.Error(nil, fmt.Sprintf(
			"failed to get token, upstream failed: code: %d, msg: %s", result.Code, result.Msg,
		))
		return errors.ErrTokenAcquisitionFailed.
			WithApiCode(result.Code).
			WithApiMessage(result.Msg).
			WithRequestID(result.RequestID).
			WithResponse(result)
	}

	// Build token struct
	var token struct {
		Token string `json:"token"`
	}

	// Unmarshal token data
	if err := result.Unmarshal(&token); err != nil {
		c.Logger.Error(nil, fmt.Sprintf(
			"failed to get token, unmarshal error: %s", err.Error(),
		))
		return errors.ErrTokenUnmarshalFailed.
			WithApiCode(result.Code).
			WithApiMessage(result.Msg).
			WithRequestID(result.RequestID).
			WithResponse(result)
	}

	// Save token
	c.token = token.Token
	return nil
}

// NewClient creates a new client to use service of Ghink Open API
func NewClient(secretID string, secretKey string, options ...Option) (*Client, error) {
	// Create client
	client := new(Client)

	// Load default logger
	client.Logger = NewLogger()

	// Load default endpoint
	client.endpoint = openapi.Endpoint

	// Load default marshal and unmarshal lib
	client.marshal = json.Marshal
	client.unmarshal = json.Unmarshal

	// Load default timeout and retry configuration (all in seconds)
	client.timeout = DefaultTimeoutSeconds                // 3 seconds timeout
	client.maxRetries = DefaultMaxRetries                 // 5 retry attempts
	client.retryDelay = DefaultRetryDelaySeconds          // 1-second initial delay
	client.exponentialBackoff = DefaultExponentialBackoff // Enable exponential backoff

	// Enable token in default
	client.enableToken = true

	// Load options
	for _, f := range options {
		f(client)
	}

	// Save keys
	client.secretID = secretID
	client.secretKey = secretKey

	// Try to get token
	if client.enableToken {
		if err := applyToken(client); err != nil {
			return nil, err
		}
	}

	return client, nil
}

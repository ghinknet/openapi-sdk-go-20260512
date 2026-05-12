package errors

import "go.gh.ink/toolbox/pointer"

// SdkError represents an SDK error with detailed information
// It contains error message, API error details, request ID, and full response
type SdkError struct {
	message string

	// API (business layer) error details
	apiCode    int    // Business error code from API response (not HTTP status code)
	apiMessage string // Business error message from API response
	requestID  string // Unique request ID for tracking

	// Full response for debugging
	response any

	raw error
}

func (e *SdkError) Error() string {
	return e.message
}

func (e *SdkError) Is(err error) bool {
	// Check if comparing against self
	if e == err {
		return true
	}
	// Check if comparing against the raw sentinel error (for cloned errors)
	return e.raw == err
}

func (e *SdkError) Unwrap() error {
	return e.raw
}

// clone creates a copy for chaining, preserving the reference to the original sentinel error.
// Internal helper used by all With* methods to ensure errors.Is always works.
func (e *SdkError) clone() *SdkError {
	ne := pointer.Copy(e)
	if ne.raw == nil {
		ne.raw = e
	}
	return ne
}

// WithApiCode sets the API error code (business layer error code, not HTTP status code)
func (e *SdkError) WithApiCode(apiCode int) *SdkError {
	ne := e.clone()
	ne.apiCode = apiCode
	return ne
}

// ApiCode returns the API error code (business layer error code, not HTTP status code)
func (e *SdkError) ApiCode() int {
	return e.apiCode
}

// WithApiMessage sets the API error message
func (e *SdkError) WithApiMessage(apiMessage string) *SdkError {
	ne := e.clone()
	ne.apiMessage = apiMessage
	return ne
}

// ApiMessage returns the API error message
func (e *SdkError) ApiMessage() string {
	return e.apiMessage
}

// WithRequestID sets the request ID for tracking
func (e *SdkError) WithRequestID(requestID string) *SdkError {
	ne := e.clone()
	ne.requestID = requestID
	return ne
}

// RequestID returns the request ID for tracking
func (e *SdkError) RequestID() string {
	return e.requestID
}

// WithResponse sets the full response data for debugging
func (e *SdkError) WithResponse(response any) *SdkError {
	ne := e.clone()
	ne.response = response
	return ne
}

// Response returns the full response data
func (e *SdkError) Response() any {
	return e.response
}

type Option func(*SdkError)

// WithApiCode option for New()
func WithApiCode(apiCode int) Option {
	return func(e *SdkError) {
		e.apiCode = apiCode
	}
}

// WithApiMessage option for New()
func WithApiMessage(apiMessage string) Option {
	return func(e *SdkError) {
		e.apiMessage = apiMessage
	}
}

// WithRequestID option for New()
func WithRequestID(requestID string) Option {
	return func(e *SdkError) {
		e.requestID = requestID
	}
}

// WithResponse option for New()
func WithResponse(response any) Option {
	return func(e *SdkError) {
		e.response = response
	}
}

func New(c string, options ...Option) *SdkError {
	err := &SdkError{message: c}

	for _, option := range options {
		option(err)
	}

	return err
}

// Predefined SDK errors for common scenarios
var (

	// Token and authentication errors

	ErrTokenAcquisitionFailed = New("failed to acquire authentication token")
	ErrTokenUnmarshalFailed   = New("failed to unmarshal token response")

	// Request and response errors

	ErrRequestSendFailed       = New("failed to send request")
	ErrResponseUnmarshalFailed = New("failed to unmarshal response body")
	ErrStatusError             = New("open api service returned status code")
)

package client

// Default timeout and retry constants (all in seconds)
const (
	// DefaultTimeoutSeconds is the default timeout for HTTP requests (in seconds)
	DefaultTimeoutSeconds = 3

	// DefaultMaxRetries is the default maximum number of retry attempts
	DefaultMaxRetries = 5

	// DefaultRetryDelaySeconds is the default initial delay between retries (in seconds)
	DefaultRetryDelaySeconds = 1

	// DefaultExponentialBackoff enables exponential backoff for retries by default
	DefaultExponentialBackoff = true
)

const (
	AuthTypeToken AuthType = iota
	AuthTypeKey
)

// MaxRetryDelaySeconds is the maximum delay between retry attempts (in seconds)
// Retry delays grow exponentially with exponential backoff but are capped at this value
const MaxRetryDelaySeconds = 60

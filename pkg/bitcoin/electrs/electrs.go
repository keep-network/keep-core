package electrs

import (
	"time"
)

const (
	// DefaultRequestTimeout is a default timeout used for HTTP requests.
	DefaultRequestTimeout = 5 * time.Second
	// DefaultRetryTimeout is a default timeout used for retries.
	DefaultRetryTimeout = 60 * time.Second
)

// Config holds configurable properties.
type Config struct {
	URL            string
	RequestTimeout time.Duration
	RetryTimeout   time.Duration
}

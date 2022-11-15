package electrum

import "time"

const (
	// DefaultKeepAliveInterval is a default interval used for Electrum server
	// connection keep alive requests.
	DefaultKeepAliveInterval = 5 * time.Minute

	// DefaultRequestRetryTimeout is a default timeout used for Electrum request
	// retries.
	DefaultRequestRetryTimeout = 60 * time.Second
)

// Config holds configurable properties.
type Config struct {
	// URL to the Electrum server in format: `hostname:port`.
	URL string
	// Electrum server protocol connection (`TCP` or `SSL`).
	Protocol Protocol
	// Interval for connection keep alive requests.
	// An Electrum server may disconnect clients that have not sent any requests
	// for roughly 10 minutes.
	KeepAliveInterval time.Duration
	// Timeout for Electrum requests retries.
	RequestRetryTimeout time.Duration
}

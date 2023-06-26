package spv

import (
	"time"
)

// Config holds configurable properties.
type Config struct {
	// Enabled indicates whether the SPV maintainer should be started.
	Enabled bool

	// RestartBackOffTime is a restart backoff which should be applied when the
	// SPV maintainer is restarted. It helps to avoid being flooded with error
	// logs in case of a permanent error in the SPV maintainer.
	RestartBackOffTime time.Duration
}

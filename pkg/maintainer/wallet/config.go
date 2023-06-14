package wallet

import "time"

const (
	DefaultRedemptionInterval = 3 * time.Hour
	DefaultSweepInterval      = 48 * time.Hour
)

// Config holds configurable properties.
type Config struct {
	Enabled            bool
	RedemptionInterval time.Duration
	SweepInterval      time.Duration
}

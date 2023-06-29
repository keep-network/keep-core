package wallet

import "time"

const (
	DefaultRedemptionInterval   = 3 * time.Hour
	DefaultDepositSweepInterval = 48 * time.Hour
)

// Config holds configurable properties.
type Config struct {
	Enabled              bool
	RedemptionInterval   time.Duration
	DepositSweepInterval time.Duration
}

package wallet

import "time"

const (
	DefaultRedemptionInterval     = 3 * time.Hour
	DefaultRedemptionWalletsLimit = 3
	DefaultDepositSweepInterval   = 48 * time.Hour
)

// Config holds configurable properties.
type Config struct {
	Enabled                bool
	RedemptionInterval     time.Duration
	RedemptionWalletsLimit uint16
	DepositSweepInterval   time.Duration
}

package wallet

import "time"

const (
	DefaultRedemptionInterval           = 3 * time.Hour
	DefaultRedemptionWalletsLimit       = 3
	DefaultRedemptionRequestAmountLimit = uint64(10 * 1e8) // 10 BTC
	DefaultDepositSweepInterval         = 48 * time.Hour
)

// Config holds configurable properties.
type Config struct {
	Enabled                      bool
	RedemptionInterval           time.Duration
	RedemptionWalletsLimit       uint16
	RedemptionRequestAmountLimit uint64
	DepositSweepInterval         time.Duration
}

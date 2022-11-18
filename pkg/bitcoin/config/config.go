package config

import (
	"github.com/keep-network/keep-core/pkg/bitcoin/electrum"
)

// Config defines the configuration for Bitcoin.
type Config struct {
	// Electrum defines the configuration for the Electrum client.
	Electrum electrum.Config
}

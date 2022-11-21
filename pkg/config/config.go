package config

import (
	"github.com/keep-network/keep-core/pkg/bitcoin/electrum"
)

// BitcoinConfig defines the configuration for Bitcoin.
type BitcoinConfig struct {
	// Electrum defines the configuration for the Electrum client.
	Electrum electrum.Config
}

package electrum

// BitcoinConfig defines the configuration for Bitcoin.
type BitcoinConfig struct {
	// Electrum defines the configuration for the Electrum client.
	Electrum Config
}

// Config holds configurable properties.
type Config struct {
	// URL to the Electrum server in format: `hostname:port`.
	URL string
}

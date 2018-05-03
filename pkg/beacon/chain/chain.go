package chain

// BeaconConfig contains configuration for the threshold relay beacon, typically
// from the underlying blockchain.
type BeaconConfig struct {
	GroupSize int
	Threshold int
}

// GetBeaconConfig Get the latest threshold relay beacon configuration.
// TODO Make this actually look up/update from chain information.
func GetBeaconConfig() BeaconConfig {
	return BeaconConfig{10, 4}
}

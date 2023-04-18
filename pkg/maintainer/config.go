package maintainer

// Config contains maintainer configuration.
type Config struct {
	// BitcoinDifficulty indicates whether the Bitcoin difficulty maintainer
	// should be started.
	BitcoinDifficulty bool

	// Wallet indicates whether the wallet maintainer should be started.
	Wallet bool
}

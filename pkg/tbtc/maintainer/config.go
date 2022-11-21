package maintainer

// Config contains maintainer configuration.
type Config struct {
	// BitcoinDifficulty indicates whether the Bitcoin difficulty maintainer
	// should be started.
	BitcoinDifficulty bool

	// TODO: Add options for other maintainer tasks, e.g. spv
}

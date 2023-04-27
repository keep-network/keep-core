package maintainer

// Config contains maintainer configuration.
type Config struct {
	// BitcoinDifficulty indicates whether the Bitcoin difficulty maintainer
	// should be started.
	BitcoinDifficulty bool

	// UseBitcoinDifficultyProxy indicates whether the Bitcoin difficulty
	// maintainer should submit difficulty information via proxy.
	UseBitcoinDifficultyProxy bool

	// TODO: Add options for other maintainer tasks, e.g. spv
}

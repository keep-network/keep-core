package maintainer

// Config contains maintainer configuration.
type Config struct {
	// Relay indicates whether maintainer should start the header relay.
	Relay bool

	// TODO: Add options for other maintainer tasks, e.g. spv
}

package tecdsa

import "time"

const (
	DefaultPreParamsPoolSize              = 50
	DefaultPreParamsGenerationTimeout     = 2 * time.Minute
	DefaultPreParamsGenerationConcurrency = 1
)

// Config carries the config for threshold ECDSA protocol.
type Config struct {
	// The size of the pre-parameters pool.
	PreParamsPoolSize int
	// Timeout for pre-parameters generation.
	PreParamsGenerationTimeout time.Duration
	// Concurrency level for pre-parameters generation.
	PreParamsGenerationConcurrency int
}

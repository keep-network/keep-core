package storage

// Config stores meta-info about keeping data on disk
type Config struct {
	// Path to the persistent storage directory on disk.
	Dir string
}

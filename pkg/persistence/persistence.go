package persistence

// Handle is an interface for data persistence. Underlying implementation
// can write data e.g. to disk, cache, or hardware module.
type Handle interface {
	Save(data []byte, directory string, name string) error
}

package persistence

// Handle is an interface for data persistence. Underlying implementation
// can write data e.g. to disk, cache, or hardware module.
type Handle interface {
	Save(data []byte, directory string, name string) error
	ReadAll(directory string) ([][]byte, error)
	Archive(from string, to string) error
	CreateDir(base string, name string) error
	GetDataDir() string
}

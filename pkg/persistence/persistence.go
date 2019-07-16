package persistence

import "github.com/ipfs/go-log"

var logger = log.Logger("keep-persistence")

// Handle is an interface for data persistence. Underlying implementation
// can write data e.g. to disk, cache, or hardware module.
type Handle interface {

	// Save takes the provided data and persists it under the given name in the
	// provided directory appropriate for the given persistent storage
	// implementation.
	Save(data []byte, directory string, name string) error

	// ReadAll returns all non-archived data. It returns two channels: the first
	// channel returned is a non-buffered channel streaming DataDescriptors of
	// all data read. The second channel is a non-buffered channel streaming all
	// errors occurred during reading. Returned channels can be integrated
	// in a pipeline pattern. The function is non-blocking. Channels are closed
	// when there is no more to be read.
	ReadAll() (<-chan DataDescriptor, <-chan error)

	// Archive marks the entire directory with the name provided as archived
	// so that the data in that directory is not returned from ReadAll.
	Archive(directory string) error
}

// DataDescriptor is an interface representing data saved in the persistence
// layer represented by Handle.
type DataDescriptor interface {
	Name() string
	Directory() string
	Content() ([]byte, error)
}

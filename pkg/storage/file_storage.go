package storage

import (
	"github.com/keep-network/keep-core/config"
)

// Storage is an interface to persist data on disk
type Storage interface {
	Save(data []byte, name string)
}

// FileStorage struct is an implementation of Storage
type FileStorage struct {
	dataDir string
}

// NewFileStorage creates a new FileStorage
func NewFileStorage() *FileStorage {
	cfg, _ := config.ReadConfig("../../../../config.local.1.toml") // path is probably different on test / prod env.

	if cfg == nil {
		return nil
	}

	return &FileStorage{
		dataDir: cfg.Storage.DataDir,
	}
}

// Save - writes data in file
func (fs *FileStorage) Save(data []byte, suffix string) {
	file := &File{
		FileName: fs.dataDir + suffix,
	}

	file.Write(data)
}

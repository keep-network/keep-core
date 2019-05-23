package storage

// Storage is an interface to persist data on disk
type Storage interface {
	Save(data []byte, name string)
}

// FileStorage struct is an implementation of Storage
type fileStorage struct {
	dataDir string
}

// NewStorage creates a new FileStorage
func NewStorage(path string) Storage {
	return &fileStorage{
		dataDir: path,
	}
}

// Save - writes data in file
func (fs *fileStorage) Save(data []byte, suffix string) {
	file := &File{
		FileName: fs.dataDir + suffix,
	}

	file.Write(data)
}

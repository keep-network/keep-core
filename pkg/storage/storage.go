package storage

// Storage is an interface to persist data on disk
type Storage interface {
	Save(data []byte, name string)
	ReadAll() [][]byte
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

// ReadAll reads all the memberships from a dir path
func (fs *fileStorage) ReadAll() [][]byte {
	file := &File{}

	return file.ReadAll(fs.dataDir)
}

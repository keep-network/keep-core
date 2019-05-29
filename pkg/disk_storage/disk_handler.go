package storage

// DiskHandler is an interface for data persistence on disk
type DiskHandler interface {
	Save(data []byte, name string) error
}

type diskStorage struct {
	dataDir string
}

// NewDiskHandler creates a new diskStorage
func NewDiskHandler(path string) DiskHandler {
	return &diskStorage{
		dataDir: path,
	}
}

// Save - writes data to file
func (ds *diskStorage) Save(data []byte, suffix string) error {
	file := &File{
		FileName: ds.dataDir + suffix,
	}

	return file.Write(data)
}

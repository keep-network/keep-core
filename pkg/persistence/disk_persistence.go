package persistence

import (
	"fmt"
	"io/ioutil"
	"os"
)

// NewDiskHandle creates on-disk data persistence handle
func NewDiskHandle(path string) Handle {
	return &diskPersistence{
		dataDir: path,
	}
}

type diskPersistence struct {
	dataDir string
}

// Save - writes data to file
func (ds *diskPersistence) Save(data []byte, suffix string) error {
	file := &file{
		fileName: ds.dataDir + suffix,
	}

	return file.Write(data)
}

var (
	//ErrNoFileExists an error is shown when no file name was provided
	errNoFileExists = fmt.Errorf("please provide a file name")
)

// File represents a file on disk that a caller can use to read and write into.
type file struct {
	// FileName is the file name of the main storage file.
	fileName string
}

// Create and write data to a file
func (f *file) Write(data []byte) error {
	if f.fileName == "" {
		return errNoFileExists
	}

	var err error
	writeFile, err := os.Create(f.fileName)
	if err != nil {
		return err
	}

	defer writeFile.Close()

	_, err = writeFile.Write(data)
	if err != nil {
		return err
	}

	writeFile.Sync()

	return nil
}

// Read a file from a file system
func (f *file) Read(fileName string) ([]byte, error) {
	if f.fileName == "" {
		return nil, errNoFileExists
	}

	readFile, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	defer readFile.Close()

	data, err := ioutil.ReadAll(readFile)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Remove a file from a file system
func (f *file) remove(fileName string) error {
	if f.fileName == "" {
		return errNoFileExists
	}

	err := os.Remove(fileName)
	if err != nil {
		return err
	}

	return nil
}

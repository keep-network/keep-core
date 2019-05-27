package storage

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	//ErrNoFileExists an error is shown when no file name was provided
	ErrNoFileExists = fmt.Errorf("please provide a file name")
)

// File represents a file on disk that a caller can use to read and write into.
type File struct {
	// FileName is the file name of the main storage file.
	FileName string
}

// NewFile creates a new file at the target location on disk
func (f *File) NewFile(FileName string) *File {

	filepath.Dir(FileName)

	return &File{
		FileName: FileName,
	}
}

// Create and write data to a file
func (f *File) Write(data []byte) error {
	if f.FileName == "" {
		return ErrNoFileExists
	}

	var err error
	writeFile, err := os.Create(f.FileName)
	check(err)

	defer writeFile.Close()

	_, err = writeFile.Write(data)
	check(err)

	writeFile.Sync()

	return nil
}

// Read a file from a file system
func (f *File) Read(FileName string) ([]byte, error) {
	if f.FileName == "" {
		return nil, ErrNoFileExists
	}

	readFile, err := os.Open(FileName)
	check(err)

	defer readFile.Close()

	data, err := ioutil.ReadAll(readFile)
	check(err)

	return data, nil
}

// Remove a file from a file syste
func (f *File) Remove(FileName string) error {
	if f.FileName == "" {
		return ErrNoFileExists
	}

	err := os.Remove(FileName)
	check(err)

	return nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

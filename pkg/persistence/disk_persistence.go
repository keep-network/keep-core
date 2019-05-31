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
func (ds *diskPersistence) Save(data []byte, dirName string, fileName string) error {
	dirPath, err := ds.createDir(dirName)
	if err != nil {
		return err
	}

	file := &file{
		filePath: fmt.Sprintf("%s%s", dirPath, fileName),
	}

	return file.write(data)
}

func (ds *diskPersistence) createDir(dirName string) (string, error) {
	dirPath := fmt.Sprintf("%s/%s", ds.dataDir, dirName)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err = os.Mkdir(dirPath, os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("error occured while creating a dir for memberships: [%v]", err)
		}
	}

	return dirPath, nil
}

var (
	//ErrNoFileExists an error is shown when no file name was provided
	errNoFileExists = fmt.Errorf("please provide a file name")
)

// File represents a file on disk that a caller can use to read and write into.
type file struct {
	// FileName is the file name of the main storage file.
	filePath string
}

// Create and write data to a file
func (f *file) write(data []byte) error {
	if f.filePath == "" {
		return errNoFileExists
	}

	var err error
	writeFile, err := os.Create(f.filePath)
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
func (f *file) read(fileName string) ([]byte, error) {
	if f.filePath == "" {
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
	if f.filePath == "" {
		return errNoFileExists
	}

	err := os.Remove(fileName)
	if err != nil {
		return err
	}

	return nil
}

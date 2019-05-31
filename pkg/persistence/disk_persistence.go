package persistence

import (
	"fmt"
	"io/ioutil"
	"os"
)

var (
	//ErrNoFileExists an error is shown when no file name was provided
	errNoFileExists = fmt.Errorf("please provide a file name")
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

// ReadAll reads all the memberships from a dir path
func (ds *diskPersistence) ReadAll() ([][]byte, error) {
	file := &file{}

	memberships, err := file.readAll(ds.dataDir)
	if err != nil {
		return nil, fmt.Errorf("error occured while reading data from disk: [%v]", err)
	}

	return memberships, nil
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

// read a file from a file system
func (f *file) read(fileName string) ([]byte, error) {
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

// readAll returns all memberships for a Keep node
func (f *file) readAll(dirPath string) ([][]byte, error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	result := [][]byte{}

	for _, file := range files {
		if file.IsDir() {
			dir, err := ioutil.ReadDir(file.Name())
			if err != nil {
				return nil, err
			}
			for _, dirFile := range dir {
				data, err := f.read(file.Name() + "/" + dirFile.Name())
				if err != nil {
					return nil, err
				}
				result = append(result, data)
			}
		}
	}

	return result, nil
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

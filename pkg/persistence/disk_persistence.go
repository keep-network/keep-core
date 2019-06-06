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

// Save - writes data to a file
func (ds *diskPersistence) Save(data []byte, dirPath, fileName string) error {
	file := &file{
		filePath: fmt.Sprintf("%s%s", dirPath, fileName),
	}

	return file.write(data)
}

// ReadAll reads all the memberships from a dir path
func (ds *diskPersistence) ReadAll(dirPath string) ([][]byte, error) {
	file := &file{}

	memberships, err := file.readAll(dirPath)
	if err != nil {
		return nil, fmt.Errorf("error occured while reading data from disk: [%v]", err)
	}

	return memberships, nil
}

// Archive a file
func (ds *diskPersistence) Archive(fromDir, toDir string) error {
	file := &file{}

	return file.archive(fromDir, toDir)
}

// CreateDir creates a directory by giving a dir path and a new dir name
func (ds *diskPersistence) CreateDir(dirBasePath, dirName string) error {
	dirPath := fmt.Sprintf("%s/%s", dirBasePath, dirName)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err = os.Mkdir(dirPath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("error occured while creating a dir: [%v]", err)
		}
	}

	return nil
}

func (ds *diskPersistence) GetDataDir() string {
	return ds.dataDir
}

// File represents a file on disk that a caller can use to write/read into or remove it.
type file struct {
	// a file that an operation is executed upon
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

func (f *file) readAll(dirPath string) ([][]byte, error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("error occured while reading dir: [%v]", err)
	}

	result := [][]byte{}

	for _, file := range files {
		if file.IsDir() {
			dir, err := ioutil.ReadDir(fmt.Sprintf("%s/%s", dirPath, file.Name()))
			if err != nil {
				return nil, fmt.Errorf("error occured while reading a directory: [%v]", err)
			}
			for _, dirFile := range dir {
				data, err := f.read(fmt.Sprintf("%s/%s/%s", dirPath, file.Name(), dirFile.Name()))
				if err != nil {
					return nil, fmt.Errorf("error occured while reading a file in directory: [%v]", err)
				}
				result = append(result, data)
			}
		}
	}

	return result, nil
}

func (f *file) archive(fromDir, toDir string) error {
	err := os.Rename(fromDir, toDir)
	if err != nil {
		return fmt.Errorf("error occured while moving a dir: [%v]", err)
	}

	return nil
}

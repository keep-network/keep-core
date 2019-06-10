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

const (
	currentDir = "current"
	archiveDir = "archive"
)

// NewDiskHandle creates on-disk data persistence handle
func NewDiskHandle(path string) Handle {
	err := createDir(path, currentDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error occured while creating [%v] directory: [%v]", currentDir, err)
	}

	err = createDir(path, archiveDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error occured while creating [%v] directory: [%v]", archiveDir, err)
	}

	return &diskPersistence{
		dataDir: path,
	}
}

type diskPersistence struct {
	dataDir string
}

// Save - writes data to the 'current' directory
func (ds *diskPersistence) Save(data []byte, dirName, fileName string) error {
	dirPath := fmt.Sprintf("%s/%s", ds.dataDir, currentDir)
	err := createDir(dirPath, dirName)
	if err != nil {
		return err
	}

	file := &file{
		filePath: fmt.Sprintf("%s/%s%s", dirPath, dirName, fileName),
	}

	return file.write(data)
}

// ReadAll data from the 'current' directory
func (ds *diskPersistence) ReadAll() ([][]byte, error) {
	file := &file{
		filePath: fmt.Sprintf("%s/%s", ds.dataDir, currentDir),
	}

	memberships, err := file.readAll()
	if err != nil {
		return nil, fmt.Errorf("error occured while reading data from disk: [%v]", err)
	}

	return memberships, nil
}

// Archive a file to the 'archive' directory
func (ds *diskPersistence) Archive(dir string) error {
	file := &file{
		filePath: ds.dataDir,
	}

	return file.archive(dir)
}

func createDir(dirBasePath, newDirName string) error {
	dirPath := fmt.Sprintf("%s/%s", dirBasePath, newDirName)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err = os.Mkdir(dirPath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("error occured while creating a dir: [%v]", err)
		}
	}

	return nil
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

func (f *file) readAll() ([][]byte, error) {
	files, err := ioutil.ReadDir(f.filePath)
	if err != nil {
		return nil, fmt.Errorf("error occured while reading dir: [%v]", err)
	}

	result := [][]byte{}

	for _, file := range files {
		if file.IsDir() {
			dir, err := ioutil.ReadDir(fmt.Sprintf("%s/%s", f.filePath, file.Name()))
			if err != nil {
				return nil, fmt.Errorf("error occured while reading a directory: [%v]", err)
			}
			for _, dirFile := range dir {
				data, err := f.read(fmt.Sprintf("%s/%s/%s", f.filePath, file.Name(), dirFile.Name()))
				if err != nil {
					return nil, fmt.Errorf("error occured while reading a file in directory: [%v]", err)
				}
				result = append(result, data)
			}
		}
	}

	return result, nil
}

func (f *file) archive(dir string) error {
	from := fmt.Sprintf("%s/%s/%s", f.filePath, currentDir, dir)
	to := fmt.Sprintf("%s/%s/%s", f.filePath, archiveDir, dir)

	err := os.Rename(from, to)
	if err != nil {
		return fmt.Errorf("error occured while moving a dir: [%v]", err)
	}

	return nil
}

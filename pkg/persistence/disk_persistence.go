package persistence

import (
	"fmt"
	"io/ioutil"
	"os"
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
	dirPath := ds.getStorageCurrentDirPath()
	err := createDir(dirPath, dirName)
	if err != nil {
		return err
	}

	return write(fmt.Sprintf("%s/%s%s", dirPath, dirName, fileName), data)
}

// ReadAll data from the 'current' directory
func (ds *diskPersistence) ReadAll() ([][]byte, error) {
	memberships, err := readAll(ds.getStorageCurrentDirPath())
	if err != nil {
		return nil, fmt.Errorf("error occured while reading data from disk: [%v]", err)
	}

	return memberships, nil
}

// Archive a directory from 'current' to 'archive'
func (ds *diskPersistence) Archive(directory string) error {
	from := fmt.Sprintf("%s/%s/%s", ds.dataDir, currentDir, directory)
	to := fmt.Sprintf("%s/%s/%s", ds.dataDir, archiveDir, directory)

	return move(from, to)
}

func (ds *diskPersistence) getStorageCurrentDirPath() string {
	return fmt.Sprintf("%s/%s", ds.dataDir, currentDir)
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

// create and write data to a file
func write(filePath string, data []byte) error {
	writeFile, err := os.Create(filePath)
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
func read(filePath string) ([]byte, error) {
	readFile, err := os.Open(filePath)
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

func readAll(directoryPath string) ([][]byte, error) {
	files, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		return nil, fmt.Errorf("error occured while reading dir: [%v]", err)
	}

	result := [][]byte{}

	for _, file := range files {
		if file.IsDir() {
			dir, err := ioutil.ReadDir(fmt.Sprintf("%s/%s", directoryPath, file.Name()))
			if err != nil {
				return nil, fmt.Errorf("error occured while reading a directory: [%v]", err)
			}
			for _, dirFile := range dir {
				data, err := read(fmt.Sprintf("%s/%s/%s", directoryPath, file.Name(), dirFile.Name()))
				if err != nil {
					return nil, fmt.Errorf("error occured while reading a file in directory: [%v]", err)
				}
				result = append(result, data)
			}
		}
	}

	return result, nil
}

func move(directoryFromPath, directoryToPath string) error {
	err := os.Rename(directoryFromPath, directoryToPath)
	if err != nil {
		return fmt.Errorf("error occured while moving a dir: [%v]", err)
	}

	return nil
}

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

func (ds *diskPersistence) Save(data []byte, dirName, fileName string) error {
	dirPath := ds.getStorageCurrentDirPath()
	err := createDir(dirPath, dirName)
	if err != nil {
		return err
	}

	return write(fmt.Sprintf("%s/%s%s", dirPath, dirName, fileName), data)
}

func (ds *diskPersistence) ReadAll() (<-chan DataDescriptor, <-chan error) {
	return readAll(ds.getStorageCurrentDirPath())
}

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

// readAll reads all files from the provided directoryPath and outputs them
// as DataDescriptors into the first returned output channel. All errors
// occurred during file system reading are sent to the second output channel
// returned from this function. The output can be later processed using
// pipeline pattern. This function is non-blocking and returned channels are
// not buffered. Channels are closed when there is no more to be read.
func readAll(directoryPath string) (<-chan DataDescriptor, <-chan error) {
	dataChannel := make(chan DataDescriptor)
	errorChannel := make(chan error)

	go func() {
		defer close(dataChannel)
		defer close(errorChannel)

		files, err := ioutil.ReadDir(directoryPath)
		if err != nil {
			errorChannel <- fmt.Errorf(
				"could not read the directory [%v]: [%v]",
				directoryPath,
				err,
			)
		}

		for _, file := range files {
			if file.IsDir() {
				dir, err := ioutil.ReadDir(fmt.Sprintf("%s/%s", directoryPath, file.Name()))
				if err != nil {
					errorChannel <- fmt.Errorf(
						"could not read the directory [%s/%s]: [%v]",
						directoryPath,
						file.Name(),
						err,
					)
				}

				for _, dirFile := range dir {
					// capture shared loop variables for the closure
					dirName := file.Name()
					fileName := dirFile.Name()

					readFunc := func() ([]byte, error) {
						return read(fmt.Sprintf(
							"%s/%s/%s",
							directoryPath,
							dirName,
							fileName,
						))
					}
					dataChannel <- &dataDescriptor{dirName, fileName, readFunc}
				}
			}
		}
	}()

	return dataChannel, errorChannel
}

func move(directoryFromPath, directoryToPath string) error {
	err := os.Rename(directoryFromPath, directoryToPath)
	if err != nil {
		return fmt.Errorf("error occured while moving a dir: [%v]", err)
	}

	return nil
}

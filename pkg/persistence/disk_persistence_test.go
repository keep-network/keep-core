package persistence

import (
	"bytes"
	"os"
	"testing"
)

var (
	fileName = "foo"
)

func TestMain(m *testing.M) {
	code := m.Run()
	if _, err := os.Stat(fileName); err == nil {
		os.Remove(fileName)
	}
	os.Exit(code)
}

func TestFile_WriteRead(t *testing.T) {
	file := &file{
		fileNamePath: fileName,
	}
	bytesToTest := []byte{115, 111, 109, 101, 10}

	file.write(bytesToTest)

	actual, _ := file.read(fileName)

	if !bytes.Equal(bytesToTest, actual) {
		t.Fatalf("Bytes do not match. \nExpected: [%+v]\nActual:   [%+v]",
			bytesToTest,
			actual)
	}
}

func TestFile_Remove(t *testing.T) {
	if _, err := os.Stat(fileName); err == nil {
		err = os.Remove(fileName)
		if err != nil {
			t.Fatalf("Was not able to remove a file [%+v]", fileName)
		}
	}

	if _, err := os.Stat(fileName); err == nil {
		t.Fatalf("File [%+v] was supposed to be removed", fileName)
	}
}

func TestDiskPersistence_Save(t *testing.T) {
	dataDir := "./"
	diskPersistence := NewDiskHandle(dataDir)
	dirName := "0x42424242"
	fileName := "/membership_test"
	pathToFile := dataDir + dirName + fileName
	bytesToTest := []byte{115, 111, 109, 101, 10}

	diskPersistence.Save(bytesToTest, dirName, fileName)

	if _, err := os.Stat(pathToFile); os.IsNotExist(err) {
		t.Fatalf("file [%+v] was supposed to be created", pathToFile)
	}

	if _, err := os.Stat(pathToFile); !os.IsNotExist(err) {
		os.RemoveAll(dataDir + dirName)
	}

	if _, err := os.Stat(pathToFile); !os.IsNotExist(err) {
		t.Fatalf("Dir [%+v] was supposed to be removed", pathToFile)
	}
}

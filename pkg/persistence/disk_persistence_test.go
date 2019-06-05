package persistence

import (
	"bytes"
	"os"
	"testing"
)

var (
	dataDir = "./"

	dirName1   = "0x424242"
	fileName11 = "/file11"
	fileName12 = "/file12"

	dirName2   = "0x777777"
	fileName21 = "/file21"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.RemoveAll(dirName1)
	os.RemoveAll(dirName2)
	os.Exit(code)
}

func TestDiskPersistence_Save(t *testing.T) {
	diskPersistence := NewDiskHandle(dataDir)
	pathToDir := dataDir + "/" + dirName1
	pathToFile := pathToDir + fileName11
	bytesToTest := []byte{115, 111, 109, 101, 10}

	diskPersistence.Save(bytesToTest, dirName1, fileName11)

	if _, err := os.Stat(pathToFile); os.IsNotExist(err) {
		t.Fatalf("file [%+v] was supposed to be created", pathToFile)
	}

	if _, err := os.Stat(pathToFile); !os.IsNotExist(err) {
		os.RemoveAll(pathToDir)
	}

	if _, err := os.Stat(pathToFile); !os.IsNotExist(err) {
		t.Fatalf("Dir [%+v] was supposed to be removed", pathToFile)
	}
}

func TestDiskPersistence_ReadAll(t *testing.T) {
	diskPersistence := NewDiskHandle(dataDir)

	bytesToTest := []byte{115, 111, 109, 101, 10}
	expectedBytes := [][]byte{bytesToTest, bytesToTest, bytesToTest}

	diskPersistence.Save(bytesToTest, dirName1, fileName11)
	diskPersistence.Save(bytesToTest, dirName1, fileName12)
	diskPersistence.Save(bytesToTest, dirName2, fileName21)

	actual, _ := diskPersistence.ReadAll()

	if len(actual) != 3 {
		t.Fatalf("Number of membership does not match. \nExpected: [%+v]\nActual:   [%+v]",
			3,
			len(actual))
	}

	for i := 0; i < 3; i++ {
		if !bytes.Equal(expectedBytes[i], actual[i]) {
			t.Fatalf("Bytes do not match. \nExpected: [%+v]\nActual:   [%+v]",
				bytesToTest,
				actual)
		}
	}
}

func TestDiskPersistence_Remove(t *testing.T) {
	diskPersistence := NewDiskHandle(dataDir)
	pathToDir := dataDir + "/" + dirName1

	bytesToTest := []byte{115, 111, 109, 101, 10}

	diskPersistence.Save(bytesToTest, dirName1, fileName11)

	if _, err := os.Stat(pathToDir); os.IsNotExist(err) {
		if err != nil {
			t.Fatalf("Dir [%+v] was supposed to be created", dirName1)
		}
	}

	diskPersistence.Remove(dirName1)

	if _, err := os.Stat(pathToDir); err == nil {
		t.Fatalf("Dir [%+v] was supposed to be removed", dirName1)
	}
}

package persistence

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

var (
	dataDir = "./"

	dirCurrent = "current"
	dirArchive = "archive"

	dirName1   = "0x424242"
	fileName11 = "/file11"
	fileName12 = "/file12"

	dirName2   = "0x777777"
	fileName21 = "/file21"

	pathToCurrent = fmt.Sprintf("%s/%s", dataDir, dirCurrent)
	pathToArchive = fmt.Sprintf("%s/%s", dataDir, dirArchive)
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.RemoveAll(pathToCurrent)
	os.RemoveAll(pathToArchive)
	os.Exit(code)
}

func TestDiskPersistence_Save(t *testing.T) {
	diskPersistence := NewDiskHandle(dataDir)

	pathToDir := fmt.Sprintf("%s/%s", pathToCurrent, dirName1)
	pathToFile := fmt.Sprintf("%s%s", pathToDir, fileName11)

	bytesToTest := []byte{115, 111, 109, 101, 10}

	diskPersistence.CreateDir(dataDir, dirCurrent)
	diskPersistence.CreateDir(pathToCurrent, dirName1)
	diskPersistence.Save(bytesToTest, pathToDir, fileName11)

	if _, err := os.Stat(pathToFile); os.IsNotExist(err) {
		t.Fatalf("file [%+v] was supposed to be created", pathToFile)
	} else {
		os.RemoveAll(pathToDir)
	}

	if _, err := os.Stat(pathToFile); !os.IsNotExist(err) {
		t.Fatalf("Dir [%+v] was supposed to be removed", pathToFile)
	}
}

func TestDiskPersistence_ReadAll(t *testing.T) {
	diskPersistence := NewDiskHandle(dataDir)

	pathToDir1 := fmt.Sprintf("%s/%s", pathToCurrent, dirName1)
	pathToDir2 := fmt.Sprintf("%s/%s", pathToCurrent, dirName2)

	bytesToTest := []byte{115, 111, 109, 101, 10}
	expectedBytes := [][]byte{bytesToTest, bytesToTest, bytesToTest}

	diskPersistence.CreateDir(dataDir, dirCurrent)
	diskPersistence.CreateDir(pathToCurrent, dirName1)
	diskPersistence.CreateDir(pathToCurrent, dirName2)

	diskPersistence.Save(bytesToTest, pathToDir1, fileName11)
	diskPersistence.Save(bytesToTest, pathToDir1, fileName12)
	diskPersistence.Save(bytesToTest, pathToDir2, fileName21)

	actual, _ := diskPersistence.ReadAll(pathToCurrent)

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

func TestDiskPersistence_Archive(t *testing.T) {
	diskPersistence := NewDiskHandle(dataDir)

	pathMoveFrom := fmt.Sprintf("%s/%s", pathToCurrent, dirName1)
	pathMoveTo := fmt.Sprintf("%s/%s", pathToArchive, dirName1)

	bytesToTest := []byte{115, 111, 109, 101, 10}

	diskPersistence.CreateDir(dataDir, dirArchive)
	diskPersistence.CreateDir(dataDir, dirCurrent)
	diskPersistence.CreateDir(pathToCurrent, dirName1)

	diskPersistence.Save(bytesToTest, pathMoveFrom, fileName11)

	if _, err := os.Stat(pathMoveFrom); os.IsNotExist(err) {
		if err != nil {
			t.Fatalf("Dir [%+v] was supposed to be created", pathMoveFrom)
		}
	}

	if _, err := os.Stat(pathMoveTo); !os.IsNotExist(err) {
		if err != nil {
			t.Fatalf("Dir [%+v] was supposed to be empty", pathMoveTo)
		}
	}

	diskPersistence.Archive(pathMoveFrom, pathMoveTo)

	if _, err := os.Stat(pathMoveFrom); !os.IsNotExist(err) {
		if err != nil {
			t.Fatalf("Dir [%+v] was supposed to be empty", pathMoveFrom)
		}
	}

	if _, err := os.Stat(pathMoveTo); os.IsNotExist(err) {
		if err != nil {
			t.Fatalf("Dir [%+v] was supposed to be created", pathMoveTo)
		}
	}
}

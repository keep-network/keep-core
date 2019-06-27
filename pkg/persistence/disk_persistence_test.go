package persistence

import (
	"fmt"
	"os"
	"testing"

	"github.com/keep-network/keep-core/pkg/internal/testutils"
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
	bytesToTest := []byte{115, 111, 109, 101, 10}

	diskPersistence.Save(bytesToTest, dirName1, fileName11)

	pathToFile := fmt.Sprintf("%s/%s%s", pathToCurrent, dirName1, fileName11)

	if _, err := os.Stat(pathToFile); os.IsNotExist(err) {
		t.Fatalf("file [%+v] was supposed to be created", pathToFile)
	}
}

func TestDiskPersistence_ReadAll(t *testing.T) {
	diskPersistence := NewDiskHandle(dataDir)

	bytesToTest := []byte{115, 111, 109, 101, 10}
	expectedBytes := [][]byte{bytesToTest, bytesToTest, bytesToTest}

	diskPersistence.Save(bytesToTest, dirName1, fileName11)
	diskPersistence.Save(bytesToTest, dirName1, fileName12)
	diskPersistence.Save(bytesToTest, dirName2, fileName21)

	dataChannel, errChannel := diskPersistence.ReadAll()

	for e := range errChannel {
		t.Error(e)
	}

	descriptors := make([]DataDescriptor, 0)
	for d := range dataChannel {
		descriptors = append(descriptors, d)
	}

	if len(descriptors) != 3 {
		t.Fatalf(
			"Number of descriptors does not match\nExpected: [%v]\nActual:   [%v]",
			3,
			len(descriptors),
		)
	}

	for i := 0; i < 3; i++ {
		fileContent, err := descriptors[i].Content()
		if err != nil {
			t.Fatal(err)
		}

		testutils.AssertBytesEqual(t, expectedBytes[i], fileContent)
	}
}

func TestDiskPersistence_Archive(t *testing.T) {
	diskPersistence := NewDiskHandle(dataDir)

	pathMoveFrom := fmt.Sprintf("%s/%s", pathToCurrent, dirName1)
	pathMoveTo := fmt.Sprintf("%s/%s", pathToArchive, dirName1)

	bytesToTest := []byte{115, 111, 109, 101, 10}

	diskPersistence.Save(bytesToTest, dirName1, fileName11)

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

	diskPersistence.Archive(dirName1)

	if _, err := os.Stat(pathMoveFrom); !os.IsNotExist(err) {
		if err != nil {
			t.Fatalf("Dir [%+v] was supposed to be moved", pathMoveFrom)
		}
	}

	if _, err := os.Stat(pathMoveTo); os.IsNotExist(err) {
		if err != nil {
			t.Fatalf("Dir [%+v] was supposed to be created", pathMoveTo)
		}
	}
}

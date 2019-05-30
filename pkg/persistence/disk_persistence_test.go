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
		fileName: fileName,
	}
	bytesToTest := []byte{115, 111, 109, 101, 10}

	file.Write(bytesToTest)

	actual, _ := file.Read(fileName)

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

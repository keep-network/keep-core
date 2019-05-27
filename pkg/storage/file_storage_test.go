package storage

import (
	"bytes"
	"os"
	"testing"
)

var (
	FileName = "foo"
)

func TestMain(m *testing.M) {
	code := m.Run()
	if _, err := os.Stat(FileName); err == nil {
		os.Remove(FileName)
	}
	os.Exit(code)
}

func TestFile_WriteRead(t *testing.T) {
	file := &File{
		FileName: FileName,
	}
	bytesToTest := []byte{115, 111, 109, 101, 10}

	file.Write(bytesToTest)

	actual, _ := file.Read(FileName)

	if !bytes.Equal(bytesToTest, actual) {
		t.Fatalf("Bytes do not match. \nExpected: [%+v]\nActual:   [%+v]",
			bytesToTest,
			actual)
	}
}

func TestFile_Remove(t *testing.T) {
	if _, err := os.Stat(FileName); err == nil {
		err = os.Remove(FileName)
		if err != nil {
			t.Fatalf("Was not able to remove a file [%+v]", FileName)
		}
	}

	if _, err := os.Stat(FileName); err == nil {
		t.Fatalf("File [%+v] was supposed to be removed", FileName)
	}
}

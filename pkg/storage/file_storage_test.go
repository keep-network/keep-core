package storage

import (
	"bytes"
	"os"
	"testing"
)

var (
	DataDir   = "./data/"
	FileName1 = DataDir + "test_foo"
	FileName2 = DataDir + "test_bar"
)

func TestMain(m *testing.M) {
	code := m.Run()
	if _, err := os.Stat(DataDir); err == nil {
		os.RemoveAll(DataDir)
	}
	os.Exit(code)
}

func TestFile_WriteRead(t *testing.T) {
	os.Mkdir(DataDir, 0777)

	file := &File{
		FileName: FileName1,
	}
	bytesToTest := []byte{115, 111, 109, 101, 10}

	file.Write(bytesToTest)

	actual, _ := file.Read(FileName1)

	if !bytes.Equal(bytesToTest, actual) {
		t.Fatalf("Bytes do not match. \nExpected: [%+v]\nActual:   [%+v]",
			bytesToTest,
			actual)
	}
}

func TestFile_Remove(t *testing.T) {
	if _, err := os.Stat(FileName1); err == nil {
		err = os.Remove(FileName1)
		if err != nil {
			t.Fatalf("Was not able to remove a file [%+v]", FileName1)
		}
	}

	if _, err := os.Stat(FileName1); err == nil {
		t.Fatalf("File [%+v] was supposed to be removed", FileName1)
	}
}

func TestFile_WriteReadAll(t *testing.T) {
	os.Mkdir(DataDir, 0777)

	file := &File{
		FileName: FileName1,
	}
	bytesToTest1 := []byte{115, 111, 109, 101, 10}
	file.Write(bytesToTest1)

	file.FileName = FileName2
	bytesToTest2 := []byte{115, 111, 109, 101, 10}
	file.Write(bytesToTest2)

	expectedResult := [][]byte{bytesToTest1, bytesToTest2}

	actual := file.ReadAll(DataDir)

	if len(expectedResult) != len(actual) {
		t.Fatalf("The size of group memberships does not match. \nExpected: [%+v]\nActual:   [%+v]",
			len(expectedResult),
			len(actual))
	}

	for i := 0; i < len(actual); i++ {
		if !bytes.Equal(expectedResult[i], actual[i]) {
			t.Fatalf("Bytes do not match. \nExpected: [%+v]\nActual:   [%+v]",
				expectedResult[i],
				actual[i])
		}
	}

}

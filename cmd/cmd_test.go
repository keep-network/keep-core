package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const (
	testFileMode = os.FileMode(0640)
)

func TestFileExists(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	file := filepath.Join(dir, "foobar")

	if FileExists(file) {
		t.Errorf("File %q should not exist", file)
	}

	err = createEmptyFile(file)
	if err != nil {
		t.Fatal(err)
	}

	if !FileExists(file) {
		t.Errorf("File %q should not exist", file)
	}
}

func TestGetConfigFilePath(t *testing.T) {
	path, err := GetConfigFilePath("")
	if err != nil {
		t.Errorf("GetConfigFilePath(\"\") returned and error: %v", err)
	} else if path != DefaultConfigPath {
		t.Errorf("Calling GetConfigFilePath(\"\"), got %s, want: %s", path, DefaultConfigPath)
	}
}

func createEmptyFile(path string) (err error) {
	return ioutil.WriteFile(path, []byte(""), testFileMode)
}

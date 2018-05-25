package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/keep-network/keep-core/pkg/conf"
)

// DefaultConfigFileName sets default file name; can be changed with --config CLI flag
const DefaultConfigFileName = "config.toml"

// GetConfigFilePath ...
func GetConfigFilePath(configPath string) (string, error) {
	if configPath == "" {
		configPath = filepath.Join(conf.CurrentDir, DefaultConfigFileName)
	}
	if exist, err := FileExists(configPath); err != nil {
		return "", err
	} else if !exist {
		return "", fmt.Errorf("config file (%s) not found", configPath)
	}
	return configPath, nil
}

// FileExists ...
func FileExists(path string) (ok bool, err error) {
	if path == "" {
		return false, errors.New("no path provided")
	}
	_, err = os.Lstat(path)
	if err != nil {
		return false, err
	}
	if os.IsNotExist(err) {
		return false, err
	}
	return true, nil
}

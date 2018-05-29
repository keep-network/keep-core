package cmd

import (
	"fmt"
	"os"
)

// GetConfigFilePath returns the full path to the project confiuration file
func GetConfigFilePath(configPath string) (string, error) {
	if configPath == "" {
		configPath = DefaultConfigPath
	}
	if exists := FileExists(configPath); !exists {
		return "", fmt.Errorf("config file (%s) not found", configPath)
	}
	return configPath, nil
}

// FileExists returns true if a file at the given path exists
func FileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

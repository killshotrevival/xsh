package utils

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func GetXSHTempDir(t *testing.T) string {
	path, err := os.MkdirTemp("", "xsh_temp_dir_*")
	if err != nil {
		t.Fatalf("error ocurred while creating temp dir for testing: %v", err)
	}
	return path
}

func RemoveTempDir(path string, t *testing.T) {
	err := os.RemoveAll(path)
	if err != nil {
		t.Fatalf("Error removing temp directory %s: %v\n", path, err)
	}
}

func ConvertToAbs(path string) (string, error) {
	// TODO: Add support for more complex relative path parsing
	if strings.Contains(path, "..") {
		return "", fmt.Errorf("`..` present in the string. Please provide absolute path")
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	path = strings.ReplaceAll(path, "~", homeDir)

	return path, nil
}

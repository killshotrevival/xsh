package utils

import (
	"os"
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

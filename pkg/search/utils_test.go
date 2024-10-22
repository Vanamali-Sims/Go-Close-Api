// utils_test.go
package search

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDirectoryExists(t *testing.T) {
	tempDir := filepath.Join(os.TempDir(), "testdir")
	os.Mkdir(tempDir, 0755)
	defer os.RemoveAll(tempDir)

	if !directoryExists(tempDir) {
		t.Errorf("directoryExists was incorrect, got: false, want: true.")
	}

	fakeDir := filepath.Join(os.TempDir(), "fakeDir")
	if directoryExists(fakeDir) {
		t.Errorf("directoryExists was incorrect, got: true, want: false.")
	}
}

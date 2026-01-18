package fs

import (
	"os"
	"path/filepath"
)

func CopyFile(from, to string) error {
	src, err := os.ReadFile(from)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(to), 0755); err != nil {
		return err
	}

	return os.WriteFile(to, src, 0644)
}

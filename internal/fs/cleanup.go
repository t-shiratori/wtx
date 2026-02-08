package fs

import (
	"os"
	"path/filepath"
)

// RemoveIfEmpty removes dir if it exists and is empty
func RemoveIfEmpty(dir string) error {
	file, err := os.Open(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	_, err = file.Readdirnames(1)
	if err == nil {
		// not empty
		return nil
	}

	// empty
	if err == os.ErrNotExist || err.Error() == "EOF" {
		return os.Remove(dir)
	}

	return err
}

// RemoveEmptyParents removes empty parent dirs up to stopDir (exclusive)
func RemoveEmptyParents(dir, stopDir string) error {
	current := filepath.Dir(dir)

	for {
		if current == stopDir || current == "." || current == "/" {
			return nil
		}

		if err := RemoveIfEmpty(current); err != nil {
			return err
		}

		current = filepath.Dir(current)
	}
}

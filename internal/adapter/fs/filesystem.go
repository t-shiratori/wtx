package fs

import (
	"os"
	"path/filepath"
)

// FileSystem implements port.FileSystem
type FileSystem struct{}

// NewFileSystem creates a new FileSystem
func NewFileSystem() *FileSystem {
	return &FileSystem{}
}

// CopyFile copies a file from src to dst
func (f *FileSystem) CopyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	return os.WriteFile(dst, data, 0644)
}

// RemoveIfEmpty removes a directory if it is empty
func (f *FileSystem) RemoveIfEmpty(dir string) error {
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

// RemoveEmptyParents removes empty parent directories up to stopDir (exclusive)
func (f *FileSystem) RemoveEmptyParents(dir, stopDir string) error {
	current := filepath.Dir(dir)

	for {
		if current == stopDir || current == "." || current == "/" {
			return nil
		}

		if err := f.RemoveIfEmpty(current); err != nil {
			return err
		}

		current = filepath.Dir(current)
	}
}

// EnsureDir creates a directory if it doesn't exist
func (f *FileSystem) EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

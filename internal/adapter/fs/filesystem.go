package fs

import (
	"os"

	oldfs "wtx/internal/fs"
)

// FileSystem implements port.FileSystem
type FileSystem struct{}

// NewFileSystem creates a new FileSystem
func NewFileSystem() *FileSystem {
	return &FileSystem{}
}

// CopyFile copies a file from src to dst
func (f *FileSystem) CopyFile(src, dst string) error {
	return oldfs.CopyFile(src, dst)
}

// RemoveIfEmpty removes a directory if it is empty
func (f *FileSystem) RemoveIfEmpty(dir string) error {
	return oldfs.RemoveIfEmpty(dir)
}

// RemoveEmptyParents removes empty parent directories up to stopDir
func (f *FileSystem) RemoveEmptyParents(dir, stopDir string) error {
	return oldfs.RemoveEmptyParents(dir, stopDir)
}

// EnsureDir creates a directory if it doesn't exist
func (f *FileSystem) EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

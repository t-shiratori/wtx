package port

// FileSystem defines the interface for file system operations
type FileSystem interface {
	// CopyFile copies a file from src to dst
	CopyFile(src, dst string) error

	// RemoveIfEmpty removes a directory if it is empty
	RemoveIfEmpty(dir string) error

	// RemoveEmptyParents removes empty parent directories up to stopDir
	RemoveEmptyParents(dir, stopDir string) error

	// EnsureDir creates a directory if it doesn't exist
	EnsureDir(path string) error
}

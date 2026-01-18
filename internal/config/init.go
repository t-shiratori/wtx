package config

import (
	"os"
	"path/filepath"
)

func EnsureWTDir(repoRoot string) error {
	return os.MkdirAll(filepath.Join(repoRoot, ".wt", "worktrees"), 0755)
}

package config

import (
	"os"
	"path/filepath"
)

const (
	DefaultWorktreesDir = "worktrees"
)

// worktreesRoot resolves the base worktrees directory
func worktreesRoot(repoRoot string, cfg *Config) string {
	root := DefaultConfigRoot
	if cfg.Worktree.RootDir != "" {
		root = cfg.Worktree.RootDir
	}
	return filepath.Join(repoRoot, root, DefaultWorktreesDir)
}

// ResolveWorktreesDir returns the resolved worktrees root directory
func ResolveWorktreesDir(repoRoot string, cfg *Config) string {
	return worktreesRoot(repoRoot, cfg)
}

// ResolveWorktreePath returns the path for a specific worktree
func ResolveWorktreePath(repoRoot string, cfg *Config, branch string) string {
	return filepath.Join(
		worktreesRoot(repoRoot, cfg),
		branch,
	)
}

// ResolveInputWorktreePath resolves a user-specified path
func ResolveInputWorktreePath(repoRoot string, cfg *Config, inputPath string) (string, error) {
	if filepath.IsAbs(inputPath) {
		return inputPath, nil
	}

	if _, err := os.Stat(inputPath); err == nil {
		return filepath.Abs(inputPath)
	}

	return filepath.Abs(
		filepath.Join(worktreesRoot(repoRoot, cfg), inputPath),
	)
}

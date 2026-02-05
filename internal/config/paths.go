package config

import (
	"os"
	"path/filepath"
)

const (
	ConfigFileName       = "config.toml"
	DefaultConfigRootDir = ".wtx"
	DefaultWorktreesDir  = "worktrees"
	DefaultWorktreeRoot  = ".wtx/worktrees"
)

// createWorktreesPath
// DefaultWorktreesDir へのパスを作成する
func createWorktreesPath(repoRoot string, rootDir string) string {
	return filepath.Join(repoRoot, rootDir, DefaultWorktreesDir)
}

// ConfigRootDir
// wtxコンフィグディレクトリのパス
func ConfigRootDir(repoRoot string) string {
	return filepath.Join(repoRoot, DefaultConfigRootDir)
}

// ConfigPath
// .wtx/config.toml
func ConfigPath(repoRoot string) string {
	return filepath.Join(ConfigRootDir(repoRoot), ConfigFileName)
}

// ResolveWorktreesDir
// worktree ルートディレクトリ（設定込みで解決）
func ResolveWorktreesDir(repoRoot string, cfg *Config) string {
	if cfg.Worktree.RootDir != "" {
		return createWorktreesPath(repoRoot, cfg.Worktree.RootDir)
	}
	return createWorktreesPath(repoRoot, DefaultConfigRootDir)
}

// ResolveWorktreePath
// 個々の worktree パス
func ResolveWorktreePath(repoRoot string, cfg *Config, branch string) string {
	return filepath.Join(
		ResolveWorktreesDir(repoRoot, cfg),
		branch,
	)
}

// ResolveInputWorktreePath
// ユーザーが指定した path を実際の worktree パスに解決する
func ResolveInputWorktreePath(repoRoot string, cfg *Config, inputPath string) (string, error) {
	if filepath.IsAbs(inputPath) {
		return inputPath, nil
	}

	if _, err := os.Stat(inputPath); err == nil {
		return filepath.Abs(inputPath)
	}

	return filepath.Abs(filepath.Join(ResolveWorktreesDir(repoRoot, cfg), inputPath))
}

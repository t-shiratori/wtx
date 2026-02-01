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

// DefaultWorktreesDir へのパスを作成する
func createWorktreesPath(repoRoot string, rootDir string) string {
	return filepath.Join(repoRoot, rootDir, DefaultWorktreesDir)
}

// wtxコンフィグディレクトリ
func ConfigRootDir(repoRoot string) string {
	return filepath.Join(repoRoot, DefaultConfigRootDir)
}

// .wtx/config.toml
func ConfigPath(repoRoot string) string {
	return filepath.Join(ConfigRootDir(repoRoot), ConfigFileName)
}

// worktree ルートディレクトリ（設定込みで解決）
func ResolveWorktreesDir(repoRoot string, cfg *Config) string {
	if cfg.Worktree.RootDir != "" {
		return createWorktreesPath(repoRoot, cfg.Worktree.RootDir)
	}
	return createWorktreesPath(repoRoot, DefaultConfigRootDir)
}

// 個々の worktree パス
func ResolveWorktreePath(repoRoot string, cfg *Config, branch string) string {
	return filepath.Join(
		ResolveWorktreesDir(repoRoot, cfg),
		branch,
	)
}

// ユーザーが指定した path を実際の worktree パスに解決する
func ResolveInputWorktreePath(repoRoot string, cfg *Config, input string) (string, error) {
	if filepath.IsAbs(input) {
		return input, nil
	}

	if _, err := os.Stat(input); err == nil {
		return filepath.Abs(input)
	}

	return filepath.Abs(filepath.Join(ResolveWorktreesDir(repoRoot, cfg), input))
}

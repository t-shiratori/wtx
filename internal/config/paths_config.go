package config

import (
	"os"
	"path/filepath"
)

const (
	AppName             = "wtx"
	GlobalConfigDirName = ".config"
	ConfigFileName      = "config.toml"
	DefaultConfigRoot   = ".wtx"
)

// localConfigDir returns the local config directory path (.wtx)
func localConfigDir(repoRoot string) string {
	return filepath.Join(repoRoot, DefaultConfigRoot)
}

// LocalConfigPath returns the local config.toml path
func LocalConfigPath(repoRoot string) string {
	return filepath.Join(localConfigDir(repoRoot), ConfigFileName)
}

// globalConfigDir returns ~/.config/wtx
func globalConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, GlobalConfigDirName, AppName), nil
}

// GlobalConfigPath returns ~/.config/wtx/config.toml
func GlobalConfigPath() (string, error) {
	dir, err := globalConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, ConfigFileName), nil
}

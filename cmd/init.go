package cmd

import (
	"os"
	"path/filepath"

	"wtx/internal/config"
	"wtx/internal/logger"

	"github.com/spf13/cobra"
)

var (
	initForce  bool
	initGlobal bool
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize wtx configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInit()
	},
}

func runInit() error {
	baseDir := config.DefaultConfigRoot

	// For global configuration
	if initGlobal {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		baseDir = filepath.Join(home, config.GlobalConfigDirName, config.AppName)
	}

	// Final config.toml path
	configPath := filepath.Join(baseDir, config.ConfigFileName)

	if _, err := os.Stat(configPath); err == nil && !initForce {
		logger.Error(os.Stdout, "%s already exists (use --force to overwrite)", configPath)
		return err
	}

	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return err
	}

	logger.Success(os.Stdout, "Initialized config at %s", configPath)

	return os.WriteFile(configPath, []byte(config.DefaultConfigTOML), 0644)
}

func init() {
	initCmd.Flags().BoolVarP(&initForce, "force", "f", false, "Overwrite existing config")
	initCmd.Flags().BoolVar(&initGlobal, "global", false, "Create global config in ~/.config/wtx")

	rootCmd.AddCommand(initCmd)
}

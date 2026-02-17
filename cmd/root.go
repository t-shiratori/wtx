package cmd

import (
	"context"
	"os"
	"wtx/internal/adapter/git"
	"wtx/internal/app"
	"wtx/internal/config"
	"wtx/internal/version"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "wtx",
	Short:   "A simple git worktree helper",
	Long:    "wtx is a small CLI tool to manage git worktrees easily.",
	Version: version.Version,
}

func init() {
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		repo := git.NewRepository()
		repoRoot, err := repo.RepoRoot()
		if err != nil {
			return err
		}

		cfg, err := config.LoadConfig(repoRoot)
		if err != nil {
			return err
		}

		worktreeDir := config.ResolveWorktreesDir(repoRoot, cfg)

		if err := config.EnsureConfigRootDir(worktreeDir); err != nil {
			return err
		}

		ctx := context.WithValue(
			cmd.Context(),
			app.Key,
			&app.Context{
				Config:   cfg,
				RepoRoot: repoRoot,
			},
		)

		cmd.SetContext(ctx)

		return nil
	}
}

// Execute is called from main.go
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

package cmd

import (
	"context"
	"os"
	"wtx/internal/adapter/git"
	"wtx/internal/app"
	"wtx/internal/config"
	"wtx/internal/shared/version"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "wtx",
	Short:   "A simple git worktree helper",
	Long:    "wtx is a small CLI tool to manage git worktrees easily.",
	Version: version.Version,
}

// init registers shared pre-run logic that runs before every subcommand (add, remove, list, etc.).
func init() {
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// Resolve the repository root; fails early if not inside a Git repo.
		repo := git.NewRepository()
		repoRoot, err := repo.RepoRoot()
		if err != nil {
			return err
		}

		// Load and merge local (.wtx/config.toml) and global (~/.config/wtx/config.toml) configs.
		// Local values take priority over global ones.
		cfg, err := config.LoadConfig(repoRoot)
		if err != nil {
			return err
		}

		// Determine the worktree placement directory from the config and repo root.
		worktreeDir := config.ResolveWorktreesDir(repoRoot, cfg)

		// Create the worktree directory if it does not exist yet.
		if err := config.EnsureConfigRootDir(worktreeDir); err != nil {
			return err
		}

		// Store app.Context (config + repo root) in Cobra's context so each subcommand
		// can retrieve it via cmd.Context().Value(app.Key).
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

package cmd

import (
	"context"
	"go-worktree-cli/internal/app"
	"go-worktree-cli/internal/config"
	"go-worktree-cli/internal/git"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "wt",
	Short: "A simple git worktree helper",
	Long:  "wt is a small CLI tool to manage git worktrees easily.",
}

func init() {
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		repoRoot, err := git.RepoRoot()
		if err != nil {
			return err
		}

		cfg, err := config.LoadConfig(repoRoot)
		if err != nil {
			return err
		}

		worktreeDir := config.ResolveWorktreesDir(repoRoot, cfg)

		if err := config.EnsureWTDir(worktreeDir); err != nil {
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

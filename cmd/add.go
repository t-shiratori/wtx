package cmd

import (
	"go-worktree-cli/internal/config"
	"go-worktree-cli/internal/fs"
	"go-worktree-cli/internal/git"
	"go-worktree-cli/internal/hook"
	"path/filepath"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <branch>",
	Short: "add git worktree",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		branch := args[0]

		// 1. git repo root
		repoRoot, err := git.RepoRoot()
		if err != nil {
			return err
		}

		// 2. load config
		cfg, err := config.LoadConfig(repoRoot)
		if err != nil {
			return err
		}

		// 3. worktree dir
		rootDir := cfg.Worktree.RootDir
		if rootDir == "" {
			rootDir = ".wt/worktrees"
		}

		worktreeDir := filepath.Join(repoRoot, rootDir, branch)

		// 4. pre hook
		if err := hook.Run(cfg.Hooks.PreCreate, repoRoot); err != nil {
			return err
		}

		// 5. git worktree add
		if err := git.AddWorktree(worktreeDir, branch); err != nil {
			return err
		}

		// 6. post create hook
		if err := hook.Run(cfg.Hooks.PostCreate, worktreeDir); err != nil {
			return err
		}

		// 7. copy files
		for _, c := range cfg.Copy {
			src := filepath.Join(repoRoot, c.From)
			dst := filepath.Join(worktreeDir, c.To)
			if err := fs.CopyFile(src, dst); err != nil {
				return err
			}
		}

		// 8. post copy hook
		if err := hook.Run(cfg.Hooks.PostCopy, worktreeDir); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

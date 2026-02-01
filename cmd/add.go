package cmd

import (
	"wtx/internal/app"
	"wtx/internal/config"
	"wtx/internal/domain/worktree"
	"wtx/internal/fs"
	"wtx/internal/git"
	"wtx/internal/hook"
	"wtx/internal/plan"
	"wtx/internal/ui"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	addDryRun    bool
	createBranch bool
	fromBranch   string
)

var addCmd = &cobra.Command{
	Use:   "add <branch>",
	Short: "add git worktree",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		appCtx := cmd.Context().Value(app.Key).(*app.Context)
		cfg := appCtx.Config
		repoRoot := appCtx.RepoRoot

		branch := args[0]

		// base branch 解決
		baseBranch := worktree.ResolveBaseBranch(fromBranch, cfg.Worktree.DefaultBaseBranch)

		// worktree dir path
		worktreeDir := config.ResolveWorktreePath(
			repoRoot,
			cfg,
			branch,
		)

		addPlan := plan.NewAddPlanFromConfig(
			cfg,
			branch,
			worktreeDir,
			createBranch,
			baseBranch,
		)

		if addDryRun {
			addPlan.Print(cmd.OutOrStdout())
			return nil
		}

		// 4. pre hook
		if err := hook.Run(cfg.Hooks.PreCreate, repoRoot); err != nil {
			return err
		}

		// 5. git worktree add
		if err := git.AddWorktree(worktreeDir, branch, createBranch, baseBranch); err != nil {
			return err
		}

		ui.Success(cmd.OutOrStdout(), "Added worktree '%s'", worktreeDir)

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
	addCmd.Flags().BoolVarP(&createBranch, "create-branch", "b", false, "create branch if it does not exist")
	addCmd.Flags().StringVar(&fromBranch, "from", "", "base branch or commit")
	addCmd.Flags().BoolVar(&addDryRun, "dry-run", false, "show what would be done")
	rootCmd.AddCommand(addCmd)
}

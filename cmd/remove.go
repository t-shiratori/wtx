package cmd

import (
	"wtx/internal/app"
	"wtx/internal/config"
	"wtx/internal/git"
	"wtx/internal/plan"
	"wtx/internal/ui"

	"github.com/spf13/cobra"
)

var (
	removeBranch bool
	force        bool
	dryRun       bool
)

var removeCmd = &cobra.Command{
	Use:   "remove <path...>",
	Short: "Remove git worktrees",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		appCtx := cmd.Context().Value(app.Key).(*app.Context)
		cfg := appCtx.Config
		repoRoot := appCtx.RepoRoot

		for _, path := range args {

			path, err := config.ResolveInputWorktreePath(repoRoot, cfg, path)
			if err != nil {
				return err
			}

			var branch string

			// 先に worktree に対応するブランチ名を取得
			if removeBranch {
				branch, err = git.BranchFromWorktree(path)
				if err != nil {
					return err
				}
			}

			plan := plan.RemovePlan{
				WorktreePath: path,
				Branch:       branch,
				Force:        force,
			}

			if dryRun {
				plan.Print(cmd.OutOrStdout())
				continue
			}

			// worktree を削除
			if err := git.RemoveWorktree(path, force); err != nil {
				return err
			}

			ui.Success(cmd.OutOrStdout(), "Removed worktree '%s'", path)

			// 対応するブランチも削除
			if removeBranch {
				if err := git.DeleteBranch(branch, force); err != nil {
					return err
				}
			}

			ui.Success(cmd.OutOrStdout(), "Deleted branch '%s'", branch)
		}
		return nil
	},
}

func init() {
	removeCmd.Flags().BoolVarP(&removeBranch, "branch", "b", false, "also delete branch")
	removeCmd.Flags().BoolVarP(&force, "force", "f", false, "force remove worktree and branch")
	removeCmd.Flags().BoolVar(&dryRun, "dry-run", false, "show what would be done")
	rootCmd.AddCommand(removeCmd)
}

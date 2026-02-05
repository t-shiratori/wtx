package cmd

import (
	"wtx/internal/app"
	"wtx/internal/config"
	"wtx/internal/git"
	"wtx/internal/plan"
	"wtx/internal/tui"
	"wtx/internal/ui"

	"github.com/spf13/cobra"
)

var (
	removeBranch bool
	force        bool
	dryRun       bool
)

var removeCmd = &cobra.Command{
	Use:   "remove [worktree ...]",
	Short: "Remove git worktrees",
	RunE: func(cmd *cobra.Command, args []string) error {

		appCtx := cmd.Context().Value(app.Key).(*app.Context)
		cfg := appCtx.Config
		repoRoot := appCtx.RepoRoot

		// --- worktree 引数が無い場合は TUI ---
		if len(args) == 0 {
			worktrees, err := git.ListWorktreesTui()
			if err != nil {
				return err
			}

			selected, err := tui.SelectWorktrees(worktrees)
			if err != nil {
				return err
			}

			if len(selected) == 0 {
				ui.Warn(cmd.OutOrStdout(), "No worktrees selected")
				return nil
			}

			args = selected
		}

		var failed int

		for _, inputPath := range args {

			err := func() error {
				// --- パス解決 ---
				path, err := config.ResolveInputWorktreePath(repoRoot, cfg, inputPath)
				if err != nil {
					return err
				}

				var branch string

				// --- ブランチ取得 ---
				if removeBranch {
					branch, err = git.BranchFromWorktree(path)
					if err != nil {
						return err
					}
				}

				// --- dryRun 用のプラン作成 ---
				plan := plan.RemovePlan{
					WorktreePath: path,
					Branch:       branch,
					Force:        force,
				}

				// --- dry-run ---
				if dryRun {
					plan.Print(cmd.OutOrStdout())
					return nil
				}

				// --- worktree 削除 ---
				if err := git.RemoveWorktree(path, force); err != nil {
					return err
				}
				ui.Success(cmd.OutOrStdout(), "Removed worktree '%s'", path)

				// --- ブランチ削除 ---
				if removeBranch && branch != "" {
					if err := git.DeleteBranch(branch, force); err != nil {
						return err
					}
					ui.Success(cmd.OutOrStdout(), "Deleted branch '%s'", branch)
				}

				return nil
			}()

			if err != nil {
				// git の stderr も含めてそのまま表示
				ui.Error(cmd.ErrOrStderr(), err.Error())
				failed++
			}
		}

		// --- まとめて警告 ---
		if failed > 0 {
			ui.Warn(
				cmd.ErrOrStderr(),
				"%d worktrees could not be removed (use --force to override)",
				failed,
			)
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

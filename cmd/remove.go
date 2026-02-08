package cmd

import (
	"wtx/internal/app"
	"wtx/internal/config"
	"wtx/internal/fs"
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

		// --- If no worktree argument is provided, use TUI ---
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
				// --- Path resolution ---
				path, err := config.ResolveInputWorktreePath(repoRoot, cfg, inputPath)
				if err != nil {
					return err
				}

				var branch string

				// --- Get branch ---
				if removeBranch {
					branch, err = git.BranchFromWorktree(path)
					if err != nil {
						return err
					}
				}

				// --- Create plan for dryRun ---
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

				// --- Remove worktree ---
				if err := git.RemoveWorktree(path, force); err != nil {
					return err
				}
				ui.Success(cmd.OutOrStdout(), "Removed worktree '%s'", path)

				// --- Delete branch ---
				if removeBranch && branch != "" {
					if err := git.DeleteBranch(branch, force); err != nil {
						return err
					}
					ui.Success(cmd.OutOrStdout(), "Deleted branch '%s'", branch)
				}

				// --- Cleanup empty directories ---
				worktreesRoot := config.ResolveWorktreesDir(repoRoot, cfg)
				_ = fs.RemoveIfEmpty(path)
				_ = fs.RemoveEmptyParents(path, worktreesRoot)

				return nil
			}()

			if err != nil {
				// Display git stderr as is
				ui.Error(cmd.ErrOrStderr(), "%v", err)
				failed++
			}
		}

		// --- Summary warning ---
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

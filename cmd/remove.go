package cmd

import (
	"fmt"
	"go-worktree-cli/internal/git"
	"go-worktree-cli/internal/ui"

	"github.com/spf13/cobra"
)

var (
	removeBranch bool
	force        bool
)

var removeCmd = &cobra.Command{
	Use:   "remove <path...>",
	Short: "Remove git worktrees",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		repoRoot, err := git.RepoRoot()
		if err != nil {
			return err
		}

		fmt.Println("Repo Root:", repoRoot)

		for _, path := range args {

			path, err := git.ResolveWorktreePath(repoRoot, path)
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

			fmt.Println("Branch:", branch)

			// worktree を削除
			if err := git.RemoveWorktree(path, force); err != nil {
				return err
			}

			ui.Success("Removed worktree '%s'", path)

			// 対応するブランチも削除
			if removeBranch {
				if err := git.DeleteBranch(branch, force); err != nil {
					return err
				}
			}

			ui.Success("Deleted branch '%s'", branch)
		}
		return nil
	},
}

func init() {
	removeCmd.Flags().BoolVar(&removeBranch, "branch", false, "also delete branch")
	removeCmd.Flags().BoolVar(&force, "force", false, "force remove worktree and branch")

	rootCmd.AddCommand(removeCmd)
}

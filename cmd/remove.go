package cmd

import (
	"go-worktree-cli/internal/git"

	"github.com/spf13/cobra"
)

var removeForce bool

var removeCmd = &cobra.Command{
	Use:   "remove <path>",
	Short: "Remove git worktree",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return git.RemoveWorktree(args[0], removeForce)
	},
}

func init() {
	removeCmd.Flags().BoolVarP(
		&removeForce,
		"force",
		"f",
		false,
		"Force remove worktree",
	)
	rootCmd.AddCommand(removeCmd)
}

package cmd

import "github.com/spf13/cobra"

var removeCmd = &cobra.Command{
	Use:   "remove <path>",
	Short: "Remove a git worktree",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return git("worktree", "remove", args[0])
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}

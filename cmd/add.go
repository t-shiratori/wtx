package cmd

import (
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <branch>",
	Short: "Add a new git worktree",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		branch := args[0]
		path := "../" + branch
		return git("worktree", "add", path, branch)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

package cmd

import (
	"wtx/internal/git"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List git worktrees",
	RunE: func(cmd *cobra.Command, args []string) error {
		return git.ListWorktrees()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

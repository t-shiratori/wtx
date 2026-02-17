package cmd

import (
	"fmt"
	"wtx/internal/adapter/git"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List git worktrees",
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := git.NewRepository()
		worktrees, err := repo.ListWorktrees()
		if err != nil {
			return err
		}

		for _, wt := range worktrees {
			fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\t%s\n", wt.Path, wt.Branch, wt.Commit)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

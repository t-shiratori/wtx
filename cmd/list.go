package cmd

import (
	"fmt"
	"text/tabwriter"
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

		// Use tabwriter to align columns with 4-space padding between them.
		w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 4, ' ', 0)
		for _, wt := range worktrees {
			fmt.Fprintf(w, "%s\t%s\t%s\n", wt.Path, wt.Branch, wt.Commit)
		}
		// Flush writes the buffered output, applying column alignment across all rows.
		w.Flush()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

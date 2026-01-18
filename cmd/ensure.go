package cmd

import (
	"fmt"
	"os"

	"go-worktree-cli/internal/git"
)

func ensureGitRepo() {
	if _, err := git.RepoRoot(); err != nil {
		fmt.Fprintln(os.Stderr, "Error: not a git repository")
		os.Exit(1)
	}
}

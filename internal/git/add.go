package git

import (
	"os"
	"os/exec"
)

// AddWorktree
// git worktree add <path> <branch>
func AddWorktree(path string, branch string) error {
	cmd := exec.Command("git", "worktree", "add", path, branch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

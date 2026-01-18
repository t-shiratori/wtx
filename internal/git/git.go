package git

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

// RepoRoot
// 現在のディレクトリが属する git リポジトリのルートを返す
func RepoRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")

	out, err := cmd.Output()
	if err != nil {
		return "", errors.New("not a git repository")
	}

	return strings.TrimSpace(string(out)), nil
}

// AddWorktree
// git worktree add <path> <branch>
func AddWorktree(path string, branch string) error {
	cmd := exec.Command("git", "worktree", "add", path, branch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func ListWorktrees() error {
	cmd := exec.Command("git", "worktree", "list")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func RemoveWorktree(path string) error {
	cmd := exec.Command("git", "worktree", "remove", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

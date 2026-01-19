package git

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ResolveWorktreePath
// ユーザーが指定した path を実際の worktree パスに解決する
func ResolveWorktreePath(repoRoot string, input string) (string, error) {
	// 1. 絶対パスならそのまま
	if filepath.IsAbs(input) {
		return input, nil
	}

	// 2. 相対パスとして実在するならそのまま
	if _, err := os.Stat(input); err == nil {
		return filepath.Abs(input)
	}

	// 3. .wt/worktrees 以下として解決
	wtPath := filepath.Join(repoRoot, ".wt", "worktrees", input)
	return filepath.Abs(wtPath)
}

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

func RemoveWorktree(path string, force bool) error {
	args := []string{"worktree", "remove"}

	if force {
		args = append(args, "--force")
	}

	args = append(args, path)

	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func DeleteBranch(branch string, force bool) error {
	flag := "-d"

	if force {
		flag = "-D"
	}

	cmd := exec.Command("git", "branch", flag, branch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func BranchFromWorktree(path string) (string, error) {

	fmt.Println("path:", path)

	cmd := exec.Command("git", "worktree", "list", "--porcelain")

	out, err := cmd.Output()

	if err != nil {
		return "", err
	}

	lines := strings.Split(string(out), "\n")

	var currentPath string
	for _, line := range lines {

		fmt.Println("line:", line)

		if rest, ok := strings.CutPrefix(line, "worktree "); ok {
			currentPath = rest
		}

		if currentPath == path {
			if rest, ok := strings.CutPrefix(line, "branch "); ok {
				branch := strings.TrimPrefix(rest, "refs/heads/")
				fmt.Println("branch:", branch)
				return branch, nil
			}
		}
	}

	fmt.Println("currentPath:", currentPath)

	return "", errors.New("branch not found for worktree")
}

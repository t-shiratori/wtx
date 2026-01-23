package git

import (
	"errors"
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
// リポジトリのルートディレクトリまでのパスを返す
func RepoRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")

	out, err := cmd.Output()
	if err != nil {
		return "", errors.New("not a git repository")
	}

	return strings.TrimSpace(string(out)), nil
}

// BranchFromWorktree
// 指定された worktree パスに対応するブランチ名を返す
func BranchFromWorktree(path string) (string, error) {

	cmd := exec.Command("git", "worktree", "list", "--porcelain")

	out, err := cmd.Output()

	if err != nil {
		return "", err
	}

	lines := strings.Split(string(out), "\n")

	var currentPath string
	for _, line := range lines {

		if rest, ok := strings.CutPrefix(line, "worktree "); ok {
			currentPath = rest
		}

		if currentPath == path {
			if rest, ok := strings.CutPrefix(line, "branch "); ok {
				branch := strings.TrimPrefix(rest, "refs/heads/")
				return branch, nil
			}
		}
	}

	return "", errors.New("branch not found for worktree")
}

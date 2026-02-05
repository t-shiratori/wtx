package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type Worktree struct {
	Path   string
	Branch string
	Commit string
}

// ListWorktreesTui
// git worktree list --porcelain を実行して worktree の一覧を取得する
func ListWorktreesTui() ([]Worktree, error) {
	cmd := exec.Command("git", "worktree", "list", "--porcelain")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git worktree list failed: %w", err)
	}

	return parsePorcelain(out)
}

// parsePorcelain
// git worktree list --porcelain の出力をパースする
func parsePorcelain(data []byte) ([]Worktree, error) {
	var result []Worktree
	var current *Worktree

	lines := bytes.Split(data, []byte("\n"))
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		fields := strings.SplitN(string(line), " ", 2)
		key := fields[0]

		switch key {
		case "worktree":
			if current != nil {
				result = append(result, *current)
			}
			current = &Worktree{
				Path: fields[1],
			}
		case "branch":
			if current != nil {
				current.Branch = strings.TrimPrefix(fields[1], "refs/heads/")
			}
		case "HEAD":
			if current != nil {
				current.Commit = fields[1]
			}
		}
	}

	if current != nil {
		result = append(result, *current)
	}

	return result, nil
}

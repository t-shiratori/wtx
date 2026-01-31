package plan

import (
	"fmt"
	"go-worktree-cli/internal/ui"
	"io"
)

type RemovePlan struct {
	WorktreePath string
	Branch       string // 空なら削除しない
	Force        bool
}

func (p *RemovePlan) Print(w io.Writer) {
	ui.Info(w, "Dry run: remove worktree")

	fmt.Fprintln(w, "  - worktree:", p.WorktreePath)

	if p.Branch != "" {
		ui.Info(w, "  - branch: %s", p.Branch)
	}

	if p.Force {
		ui.Info(w, "  - force: true")
	}
}

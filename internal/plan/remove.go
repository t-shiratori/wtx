package plan

import (
	"fmt"
	"go-worktree-cli/internal/ui"
)

type RemovePlan struct {
	WorktreePath string
	Branch       string // 空なら削除しない
	Force        bool
}

func (p *RemovePlan) Print() {
	ui.Info("Dry run: remove worktree")

	fmt.Println("  - worktree:", p.WorktreePath)

	if p.Branch != "" {
		fmt.Println("  - branch:", p.Branch)
	}

	if p.Force {
		fmt.Println("  - force: true")
	}
}

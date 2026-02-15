package plan

import (
	"io"
	"wtx/internal/infra/logger"
)

type RemovePlan struct {
	WorktreePath string
	Branch       string // 空なら削除しない
	Force        bool
}

func (p *RemovePlan) Print(w io.Writer) {
	logger.Info(w, "Dry run: remove worktree")

	logger.Info(w, "  - worktree: %s", p.WorktreePath)

	if p.Branch != "" {
		logger.Info(w, "  - branch: %s", p.Branch)
	}

	if p.Force {
		logger.Info(w, "  - force: true")
	}
}

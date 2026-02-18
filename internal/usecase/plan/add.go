package plan

import (
	"io"
	"wtx/internal/shared/logger"
)

type AddPlan struct {
	Branch       string
	BaseBranch   string
	WorktreePath string

	CreateBranch bool

	PreHook      []string
	PostHook     []string
	CopyFiles    []string
	PostCopyHook []string
}

func (p *AddPlan) Print(w io.Writer) {
	logger.Info(w, "Dry run: add worktree")

	logger.Info(w, "  - branch: %s", p.Branch)
	logger.Info(w, "  - worktree: %s", p.WorktreePath)
	logger.Info(w, "  - from: %s", p.BaseBranch)

	if p.CreateBranch {
		logger.Info(w, "  - create-branch: %s", p.Branch)
	}

	logger.Info(w, "  - hooks:")

	if len(p.PreHook) > 0 {
		logger.Info(w, "    - pre-hook:")
		for _, h := range p.PreHook {
			logger.Info(w, "      - %s", h)
		}
	}

	if len(p.PostHook) > 0 {
		logger.Info(w, "    - post-hook:")
		for _, h := range p.PostHook {
			logger.Info(w, "      - %s", h)
		}
	}

	if len(p.CopyFiles) > 0 {
		logger.Info(w, "    - copy:")
		for _, f := range p.CopyFiles {
			logger.Info(w, "      - %s", f)
		}
	}

	if len(p.PostCopyHook) > 0 {
		logger.Info(w, "    - post-copy-hook:")
		for _, h := range p.PostCopyHook {
			logger.Info(w, "      - %s", h)
		}
	}
}

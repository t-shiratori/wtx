package plan

import "wtx/internal/config"

func NewAddPlanFromConfig(
	cfg *config.Config,
	branch string,
	worktreePath string,
	createBranch bool,
	baseBranch string,
) *AddPlan {

	plan := &AddPlan{
		Branch:       branch,
		WorktreePath: worktreePath,
		CreateBranch: createBranch,
		BaseBranch:   baseBranch,
	}

	// --- hooks ---
	plan.PreHook = cfg.Hooks.PreCreate
	plan.PostHook = cfg.Hooks.PostCreate
	plan.PostCopyHook = cfg.Hooks.PostCopy

	// --- copy ---
	for _, c := range cfg.Copy {
		// 表示用なので from -> to を 1 行にまとめる
		plan.CopyFiles = append(
			plan.CopyFiles,
			c.From+" -> "+c.To,
		)
	}

	return plan
}

package worktree

// AddOptions contains parameters for adding a worktree
type AddOptions struct {
	Branch       string
	Path         string
	BaseBranch   string
	CreateBranch bool
}

// RemoveOptions contains parameters for removing a worktree
type RemoveOptions struct {
	Path         string
	DeleteBranch bool
	Force        bool
}

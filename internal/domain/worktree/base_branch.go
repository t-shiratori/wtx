package worktree

// ResolveBaseBranch determines the base branch
func ResolveBaseBranch(fromFlag string, defaultFromConfig string) string {
	if fromFlag != "" {
		return fromFlag
	}
	if defaultFromConfig != "" {
		return defaultFromConfig
	}
	return "HEAD"
}

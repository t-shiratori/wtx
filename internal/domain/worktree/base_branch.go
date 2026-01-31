package worktree

// ResolveBaseBranch は base branch を決定する
func ResolveBaseBranch(fromFlag string, defaultFromConfig string) string {
	if fromFlag != "" {
		return fromFlag
	}
	if defaultFromConfig != "" {
		return defaultFromConfig
	}
	return "HEAD"
}

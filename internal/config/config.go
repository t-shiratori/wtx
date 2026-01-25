package config

type Config struct {
	Worktree WorktreeConfig `toml:"worktree"`
	Copy     []CopyConfig   `toml:"copy"`
	Hooks    HookConfig     `toml:"hooks"`
}

type WorktreeConfig struct {
	RootDir           string `toml:"root_dir"`
	DefaultBaseBranch string `toml:"default_base_branch"`
}

type CopyConfig struct {
	From string `toml:"from"`
	To   string `toml:"to"`
}

type HookConfig struct {
	PreCreate  []string `toml:"pre_create"`
	PostCreate []string `toml:"post_create"`
	PostCopy   []string `toml:"post_copy"`
}

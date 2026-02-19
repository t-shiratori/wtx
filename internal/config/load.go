package config

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

func loadConfigFile(path string) (*Config, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // return nil if file does not exist
		}
		return nil, err
	}

	var cfg Config
	if err := toml.Unmarshal(src, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func LoadConfig(repoRoot string) (*Config, error) {
	globalPath, err := GlobalConfigPath()
	if err != nil {
		return nil, err
	}
	globalCfg, err := loadConfigFile(globalPath)
	if err != nil {
		return nil, err
	}

	localPath := LocalConfigPath(repoRoot)
	localCfg, err := loadConfigFile(localPath)
	if err != nil {
		return nil, err
	}

	return mergeConfig(globalCfg, localCfg), nil
}

// mergeConfig merges global and local configs.
// local values take priority over global values.
// Either argument can be nil.
func mergeConfig(global, local *Config) *Config {
	result := &Config{}

	if global != nil {
		result.Worktree.RootDir = global.Worktree.RootDir
		result.Worktree.DefaultBaseBranch = global.Worktree.DefaultBaseBranch
		result.Copy.Patterns = global.Copy.Patterns
		result.Copy.Files = global.Copy.Files
		result.Hooks.PreCreate = global.Hooks.PreCreate
		result.Hooks.PostCreate = global.Hooks.PostCreate
		result.Hooks.PostCopy = global.Hooks.PostCopy
	}

	if local != nil {
		if local.Worktree.RootDir != "" {
			result.Worktree.RootDir = local.Worktree.RootDir
		}
		if local.Worktree.DefaultBaseBranch != "" {
			result.Worktree.DefaultBaseBranch = local.Worktree.DefaultBaseBranch
		}
		result.Copy.Patterns = mergeStringSlices(result.Copy.Patterns, local.Copy.Patterns)
		result.Copy.Files = mergeCopyFiles(result.Copy.Files, local.Copy.Files)
		result.Hooks.PreCreate = mergeStringSlices(result.Hooks.PreCreate, local.Hooks.PreCreate)
		result.Hooks.PostCreate = mergeStringSlices(result.Hooks.PostCreate, local.Hooks.PostCreate)
		result.Hooks.PostCopy = mergeStringSlices(result.Hooks.PostCopy, local.Hooks.PostCopy)
	}

	return result
}

// mergeStringSlices combines two string slices, removing duplicates.
// Items appearing in both slices are kept once (from global position).
// Items only in local are appended after global items.
func mergeStringSlices(global, local []string) []string {
	seen := map[string]bool{}
	result := []string{}
	for _, v := range global {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	for _, v := range local {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}

// mergeCopyFiles deduplicates by From field; local entries win over global.
// Order: global entries first (replaced by local if From matches), then local-only entries.
func mergeCopyFiles(global, local []CopyFile) []CopyFile {
	localByFrom := map[string]CopyFile{}
	for _, f := range local {
		localByFrom[f.From] = f
	}

	result := []CopyFile{}
	globalFromSet := map[string]bool{}
	for _, f := range global {
		globalFromSet[f.From] = true
		if lf, ok := localByFrom[f.From]; ok {
			result = append(result, lf) // local wins
		} else {
			result = append(result, f)
		}
	}
	for _, f := range local {
		if !globalFromSet[f.From] {
			result = append(result, f)
		}
	}
	return result
}

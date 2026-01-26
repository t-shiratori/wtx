package config

import (
	"path/filepath"
	"testing"
)

func TestResolveWorktreesDir(t *testing.T) {
	repoRoot := "/repo"

	tests := []struct {
		name     string
		rootDir  string
		expected string
	}{
		{
			name:     "default root_dir",
			rootDir:  "",
			expected: filepath.Join(repoRoot, DefaultWorktreeRoot),
		},
		{
			name:     "custom root_dir",
			rootDir:  "custom/",
			expected: filepath.Join(repoRoot, "custom/worktrees"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Worktree: WorktreeConfig{
					RootDir: tt.rootDir,
				},
			}

			got := ResolveWorktreesDir(repoRoot, cfg)

			if got != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

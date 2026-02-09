package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveWorktreesDir_Default(t *testing.T) {
	repoRoot := "/repo"

	cfg := &Config{}

	got := ResolveWorktreesDir(repoRoot, cfg)
	want := filepath.Join(repoRoot, DefaultConfigRoot, DefaultWorktreesDir)

	if got != want {
		t.Fatalf("ResolveWorktreesDir() = %q, want %q", got, want)
	}
}

func TestResolveWorktreesDir_CustomRoot(t *testing.T) {
	repoRoot := "/repo"

	cfg := &Config{}
	cfg.Worktree.RootDir = ".custom"

	got := ResolveWorktreesDir(repoRoot, cfg)
	want := filepath.Join(repoRoot, ".custom", DefaultWorktreesDir)

	if got != want {
		t.Fatalf("ResolveWorktreesDir() = %q, want %q", got, want)
	}
}

func TestResolveWorktreePath(t *testing.T) {
	repoRoot := "/repo"
	branch := "feature-1"

	cfg := &Config{}

	got := ResolveWorktreePath(repoRoot, cfg, branch)
	want := filepath.Join(
		repoRoot,
		DefaultConfigRoot,
		DefaultWorktreesDir,
		branch,
	)

	if got != want {
		t.Fatalf("ResolveWorktreePath() = %q, want %q", got, want)
	}
}

func TestResolveInputWorktreePath_Absolute(t *testing.T) {
	repoRoot := "/repo"
	absPath := "/abs/worktree"

	cfg := &Config{}

	got, err := ResolveInputWorktreePath(repoRoot, cfg, absPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != absPath {
		t.Fatalf("ResolveInputWorktreePath() = %q, want %q", got, absPath)
	}
}

func TestResolveInputWorktreePath_ExistingRelativePath(t *testing.T) {
	repoRoot := t.TempDir()

	// Create existing directory
	existing := filepath.Join(repoRoot, "existing")
	if err := os.Mkdir(existing, 0o755); err != nil {
		t.Fatal(err)
	}

	cfg := &Config{}

	got, err := ResolveInputWorktreePath(repoRoot, cfg, existing)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want, _ := filepath.Abs(existing)

	if got != want {
		t.Fatalf("ResolveInputWorktreePath() = %q, want %q", got, want)
	}
}

func TestResolveInputWorktreePath_NonExistingRelative(t *testing.T) {
	repoRoot := t.TempDir()

	cfg := &Config{}
	input := "feature-2"

	got, err := ResolveInputWorktreePath(repoRoot, cfg, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want, _ := filepath.Abs(
		filepath.Join(
			repoRoot,
			DefaultConfigRoot,
			DefaultWorktreesDir,
			input,
		),
	)

	if got != want {
		t.Fatalf("ResolveInputWorktreePath() = %q, want %q", got, want)
	}
}

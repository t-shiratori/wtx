package config

import (
	"os"
	"path/filepath"
	"testing"
)

// --- mergeStringSlices ---

func TestMergeStringSlices_BothEmpty(t *testing.T) {
	got := mergeStringSlices(nil, nil)
	if len(got) != 0 {
		t.Fatalf("expected empty, got %v", got)
	}
}

func TestMergeStringSlices_GlobalOnly(t *testing.T) {
	got := mergeStringSlices([]string{"a", "b"}, nil)
	want := []string{"a", "b"}
	if !equalStringSlice(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestMergeStringSlices_LocalOnly(t *testing.T) {
	got := mergeStringSlices(nil, []string{"a", "b"})
	want := []string{"a", "b"}
	if !equalStringSlice(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestMergeStringSlices_NoDuplicates(t *testing.T) {
	got := mergeStringSlices([]string{"a", "b"}, []string{"c", "d"})
	want := []string{"a", "b", "c", "d"}
	if !equalStringSlice(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestMergeStringSlices_WithDuplicates(t *testing.T) {
	// "b" appears in both: kept once at global position
	got := mergeStringSlices([]string{"a", "b"}, []string{"b", "c"})
	want := []string{"a", "b", "c"}
	if !equalStringSlice(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

// --- mergeCopyFiles ---

func TestMergeCopyFiles_BothEmpty(t *testing.T) {
	got := mergeCopyFiles(nil, nil)
	if len(got) != 0 {
		t.Fatalf("expected empty, got %v", got)
	}
}

func TestMergeCopyFiles_GlobalOnly(t *testing.T) {
	global := []CopyFile{{From: "a", To: "a2"}}
	got := mergeCopyFiles(global, nil)
	if len(got) != 1 || got[0] != global[0] {
		t.Fatalf("got %v, want %v", got, global)
	}
}

func TestMergeCopyFiles_LocalOnly(t *testing.T) {
	local := []CopyFile{{From: "a", To: "a2"}}
	got := mergeCopyFiles(nil, local)
	if len(got) != 1 || got[0] != local[0] {
		t.Fatalf("got %v, want %v", got, local)
	}
}

func TestMergeCopyFiles_LocalWinsOnDuplicate(t *testing.T) {
	global := []CopyFile{{From: "a", To: "a_global"}}
	local := []CopyFile{{From: "a", To: "a_local"}}
	got := mergeCopyFiles(global, local)
	if len(got) != 1 || got[0].To != "a_local" {
		t.Fatalf("got %v, want To=a_local", got)
	}
}

func TestMergeCopyFiles_MergeNoDuplicate(t *testing.T) {
	global := []CopyFile{{From: "a", To: "a2"}}
	local := []CopyFile{{From: "b", To: "b2"}}
	got := mergeCopyFiles(global, local)
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %v", got)
	}
	if got[0].From != "a" || got[1].From != "b" {
		t.Fatalf("unexpected order: %v", got)
	}
}

// --- mergeConfig ---

func TestMergeConfig_BothNil(t *testing.T) {
	got := mergeConfig(nil, nil)
	if got == nil {
		t.Fatal("expected non-nil Config")
	}
	if got.Worktree.RootDir != "" || got.Worktree.DefaultBaseBranch != "" {
		t.Fatalf("expected zero Config, got %+v", got)
	}
	if len(got.Copy.Patterns) != 0 || len(got.Copy.Files) != 0 {
		t.Fatalf("expected zero Config, got %+v", got)
	}
	if len(got.Hooks.PreCreate) != 0 || len(got.Hooks.PostCreate) != 0 || len(got.Hooks.PostCopy) != 0 {
		t.Fatalf("expected zero Config, got %+v", got)
	}
}

func TestMergeConfig_GlobalOnly(t *testing.T) {
	global := &Config{}
	global.Worktree.RootDir = "global-root"
	global.Worktree.DefaultBaseBranch = "main"
	global.Hooks.PreCreate = []string{"global-pre"}

	got := mergeConfig(global, nil)

	if got.Worktree.RootDir != "global-root" {
		t.Errorf("RootDir: got %q, want %q", got.Worktree.RootDir, "global-root")
	}
	if got.Worktree.DefaultBaseBranch != "main" {
		t.Errorf("DefaultBaseBranch: got %q, want %q", got.Worktree.DefaultBaseBranch, "main")
	}
	if !equalStringSlice(got.Hooks.PreCreate, []string{"global-pre"}) {
		t.Errorf("PreCreate: got %v", got.Hooks.PreCreate)
	}
}

func TestMergeConfig_LocalOnly(t *testing.T) {
	local := &Config{}
	local.Worktree.RootDir = "local-root"
	local.Hooks.PostCreate = []string{"local-post"}

	got := mergeConfig(nil, local)

	if got.Worktree.RootDir != "local-root" {
		t.Errorf("RootDir: got %q, want %q", got.Worktree.RootDir, "local-root")
	}
	if !equalStringSlice(got.Hooks.PostCreate, []string{"local-post"}) {
		t.Errorf("PostCreate: got %v", got.Hooks.PostCreate)
	}
}

func TestMergeConfig_LocalStringOverridesGlobal(t *testing.T) {
	global := &Config{}
	global.Worktree.RootDir = "global-root"
	global.Worktree.DefaultBaseBranch = "main"

	local := &Config{}
	local.Worktree.RootDir = "local-root"
	// DefaultBaseBranch is empty in local → should fall back to global

	got := mergeConfig(global, local)

	if got.Worktree.RootDir != "local-root" {
		t.Errorf("RootDir: got %q, want %q", got.Worktree.RootDir, "local-root")
	}
	if got.Worktree.DefaultBaseBranch != "main" {
		t.Errorf("DefaultBaseBranch: got %q, want %q", got.Worktree.DefaultBaseBranch, "main")
	}
}

func TestMergeConfig_SlicesMergedWithDedup(t *testing.T) {
	global := &Config{}
	global.Hooks.PreCreate = []string{"setup", "lint"}
	global.Copy.Patterns = []string{"*.env"}

	local := &Config{}
	local.Hooks.PreCreate = []string{"lint", "test"} // "lint" is duplicate
	local.Copy.Patterns = []string{"*.env", "*.local"} // "*.env" is duplicate

	got := mergeConfig(global, local)

	wantPre := []string{"setup", "lint", "test"}
	if !equalStringSlice(got.Hooks.PreCreate, wantPre) {
		t.Errorf("PreCreate: got %v, want %v", got.Hooks.PreCreate, wantPre)
	}

	wantPatterns := []string{"*.env", "*.local"}
	if !equalStringSlice(got.Copy.Patterns, wantPatterns) {
		t.Errorf("Patterns: got %v, want %v", got.Copy.Patterns, wantPatterns)
	}
}

// --- LoadConfig integration ---

func TestLoadConfig_NeitherExists(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	repoRoot := t.TempDir()
	got, err := LoadConfig(repoRoot)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil {
		t.Fatal("expected non-nil Config")
	}
}

func TestLoadConfig_LocalOnly(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	repoRoot := t.TempDir()
	localDir := filepath.Join(repoRoot, DefaultConfigRoot)
	if err := os.MkdirAll(localDir, 0o755); err != nil {
		t.Fatal(err)
	}
	localCfgPath := filepath.Join(localDir, ConfigFileName)
	content := "[worktree]\nroot_dir = \"local-root\"\n"
	if err := os.WriteFile(localCfgPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := LoadConfig(repoRoot)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Worktree.RootDir != "local-root" {
		t.Errorf("RootDir: got %q, want %q", got.Worktree.RootDir, "local-root")
	}
}

func TestLoadConfig_MergesGlobalAndLocal(t *testing.T) {
	// Setup fake HOME to isolate global config
	fakeHome := t.TempDir()
	t.Setenv("HOME", fakeHome)

	// Write global config
	globalDir := filepath.Join(fakeHome, GlobalConfigDirName, AppName)
	if err := os.MkdirAll(globalDir, 0o755); err != nil {
		t.Fatal(err)
	}
	globalContent := "[worktree]\ndefault_base_branch = \"main\"\n\n[hooks]\npre_create = [\"global-hook\"]\n"
	if err := os.WriteFile(filepath.Join(globalDir, ConfigFileName), []byte(globalContent), 0o644); err != nil {
		t.Fatal(err)
	}

	// Write local config
	repoRoot := t.TempDir()
	localDir := filepath.Join(repoRoot, DefaultConfigRoot)
	if err := os.MkdirAll(localDir, 0o755); err != nil {
		t.Fatal(err)
	}
	localContent := "[worktree]\nroot_dir = \"local-root\"\n\n[hooks]\npre_create = [\"local-hook\"]\n"
	if err := os.WriteFile(filepath.Join(localDir, ConfigFileName), []byte(localContent), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := LoadConfig(repoRoot)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// local overrides global for RootDir
	if got.Worktree.RootDir != "local-root" {
		t.Errorf("RootDir: got %q, want %q", got.Worktree.RootDir, "local-root")
	}
	// global value used for DefaultBaseBranch (local doesn't set it)
	if got.Worktree.DefaultBaseBranch != "main" {
		t.Errorf("DefaultBaseBranch: got %q, want %q", got.Worktree.DefaultBaseBranch, "main")
	}
	// slices merged
	wantHooks := []string{"global-hook", "local-hook"}
	if !equalStringSlice(got.Hooks.PreCreate, wantHooks) {
		t.Errorf("PreCreate: got %v, want %v", got.Hooks.PreCreate, wantHooks)
	}
}

// --- helpers ---

func equalStringSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

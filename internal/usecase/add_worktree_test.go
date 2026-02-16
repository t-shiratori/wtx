package usecase_test

import (
	"testing"

	"wtx/internal/config"
	"wtx/internal/domain/worktree"
	"wtx/internal/usecase"
)

// mockGitRepo is a mock implementation of port.GitRepository
type mockGitRepo struct {
	addWorktreeCalled bool
	addWorktreeOpts   worktree.AddOptions
	addWorktreeErr    error
}

func (m *mockGitRepo) AddWorktree(opts worktree.AddOptions) error {
	m.addWorktreeCalled = true
	m.addWorktreeOpts = opts
	return m.addWorktreeErr
}

func (m *mockGitRepo) RemoveWorktree(path string, force bool) error { return nil }
func (m *mockGitRepo) DeleteBranch(branch string, force bool) error { return nil }
func (m *mockGitRepo) ListWorktrees() ([]worktree.Worktree, error)  { return nil, nil }
func (m *mockGitRepo) BranchFromWorktree(path string) (string, error) {
	return "", nil
}
func (m *mockGitRepo) RepoRoot() (string, error) { return "", nil }

// mockFS is a mock implementation of port.FileSystem
type mockFS struct {
	copyFileCalled bool
	copiedFiles    []struct{ src, dst string }
}

func (m *mockFS) CopyFile(src, dst string) error {
	m.copyFileCalled = true
	m.copiedFiles = append(m.copiedFiles, struct{ src, dst string }{src, dst})
	return nil
}
func (m *mockFS) RemoveIfEmpty(dir string) error              { return nil }
func (m *mockFS) RemoveEmptyParents(dir, stopDir string) error { return nil }
func (m *mockFS) EnsureDir(path string) error                 { return nil }

// mockHookRunner is a mock implementation of port.HookRunner
type mockHookRunner struct {
	runCalled   bool
	runCommands []string
}

func (m *mockHookRunner) Run(commands []string, dir string) error {
	m.runCalled = true
	m.runCommands = commands
	return nil
}

func TestAddWorktree_Execute(t *testing.T) {
	gitRepo := &mockGitRepo{}
	fs := &mockFS{}
	hook := &mockHookRunner{}
	cfg := &config.Config{
		Worktree: config.WorktreeConfig{
			DefaultBaseBranch: "main",
		},
	}

	uc := usecase.NewAddWorktree(gitRepo, fs, hook, cfg, "/tmp/repo")

	output, err := uc.Execute(usecase.AddWorktreeInput{
		Branch:       "feature/test",
		CreateBranch: true,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !gitRepo.addWorktreeCalled {
		t.Fatal("AddWorktree was not called")
	}

	if output.Branch != "feature/test" {
		t.Errorf("expected branch 'feature/test', got '%s'", output.Branch)
	}

	if output.BaseBranch != "main" {
		t.Errorf("expected base branch 'main', got '%s'", output.BaseBranch)
	}
}

func TestAddWorktree_FromBranchOverridesConfig(t *testing.T) {
	gitRepo := &mockGitRepo{}
	fs := &mockFS{}
	hook := &mockHookRunner{}
	cfg := &config.Config{
		Worktree: config.WorktreeConfig{
			DefaultBaseBranch: "main",
		},
	}

	uc := usecase.NewAddWorktree(gitRepo, fs, hook, cfg, "/tmp/repo")

	output, err := uc.Execute(usecase.AddWorktreeInput{
		Branch:     "feature/test",
		FromBranch: "develop",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output.BaseBranch != "develop" {
		t.Errorf("expected base branch 'develop', got '%s'", output.BaseBranch)
	}

	if gitRepo.addWorktreeOpts.BaseBranch != "develop" {
		t.Errorf("expected git opts base branch 'develop', got '%s'", gitRepo.addWorktreeOpts.BaseBranch)
	}
}

func TestAddWorktree_CopyPatternsNewFormat(t *testing.T) {
	gitRepo := &mockGitRepo{}
	fs := &mockFS{}
	hook := &mockHookRunner{}
	cfg := &config.Config{
		Worktree: config.WorktreeConfig{
			DefaultBaseBranch: "main",
		},
		Copy: config.CopyConfigSection{
			Patterns: []string{"*.env", "config/*.yaml"},
		},
	}

	uc := usecase.NewAddWorktree(gitRepo, fs, hook, cfg, "/tmp/repo")

	_, err := uc.Execute(usecase.AddWorktreeInput{
		Branch:       "feature/test",
		CreateBranch: true,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Note: Since we're using glob patterns, actual file matching depends on filesystem
	// In a real scenario with test fixtures, we'd verify files were copied
}

func TestAddWorktree_BackwardCompatibility(t *testing.T) {
	gitRepo := &mockGitRepo{}
	fs := &mockFS{}
	hook := &mockHookRunner{}
	cfg := &config.Config{
		Worktree: config.WorktreeConfig{
			DefaultBaseBranch: "main",
		},
		Copy: config.CopyConfigSection{
			Files: []config.CopyFile{
				{From: ".env.example", To: ".env"},
				{From: "config.yaml", To: "local.yaml"},
			},
		},
	}

	uc := usecase.NewAddWorktree(gitRepo, fs, hook, cfg, "/tmp/repo")

	_, err := uc.Execute(usecase.AddWorktreeInput{
		Branch:       "feature/test",
		CreateBranch: true,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !fs.copyFileCalled {
		t.Fatal("CopyFile was not called")
	}

	if len(fs.copiedFiles) != 2 {
		t.Fatalf("expected 2 files copied, got %d", len(fs.copiedFiles))
	}

	// Verify first copy
	if fs.copiedFiles[0].src != "/tmp/repo/.env.example" {
		t.Errorf("expected src '/tmp/repo/.env.example', got '%s'", fs.copiedFiles[0].src)
	}

	// Verify second copy
	if fs.copiedFiles[1].src != "/tmp/repo/config.yaml" {
		t.Errorf("expected src '/tmp/repo/config.yaml', got '%s'", fs.copiedFiles[1].src)
	}
}

func TestAddWorktree_MixedFormat(t *testing.T) {
	gitRepo := &mockGitRepo{}
	fs := &mockFS{}
	hook := &mockHookRunner{}
	cfg := &config.Config{
		Worktree: config.WorktreeConfig{
			DefaultBaseBranch: "main",
		},
		Copy: config.CopyConfigSection{
			Patterns: []string{"*.env"},
			Files: []config.CopyFile{
				{From: ".env.example", To: ".env"},
			},
		},
	}

	uc := usecase.NewAddWorktree(gitRepo, fs, hook, cfg, "/tmp/repo")

	_, err := uc.Execute(usecase.AddWorktreeInput{
		Branch:       "feature/test",
		CreateBranch: true,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Both old and new formats should be processed
	// Actual verification depends on filesystem state for glob patterns
}

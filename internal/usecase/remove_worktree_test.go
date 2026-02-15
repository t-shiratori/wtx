package usecase_test

import (
	"testing"

	"wtx/internal/config"
	"wtx/internal/usecase"
)

func TestRemoveWorktree_Execute(t *testing.T) {
	gitRepo := &mockGitRepo{}
	fs := &mockFS{}
	cfg := &config.Config{}

	uc := usecase.NewRemoveWorktree(gitRepo, fs, cfg, "/tmp/repo")

	output, err := uc.Execute(usecase.RemoveWorktreeInput{
		Paths: []string{"feature/test"},
		Force: false,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(output.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(output.Results))
	}

	if !output.Results[0].Success {
		t.Errorf("expected success, got error: %v", output.Results[0].Error)
	}
}

func TestRemoveWorktree_WithDeleteBranch(t *testing.T) {
	gitRepo := &mockGitRepo{}
	fs := &mockFS{}
	cfg := &config.Config{}

	uc := usecase.NewRemoveWorktree(gitRepo, fs, cfg, "/tmp/repo")

	output, err := uc.Execute(usecase.RemoveWorktreeInput{
		Paths:        []string{"feature/test"},
		DeleteBranch: true,
		Force:        false,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(output.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(output.Results))
	}
}

func TestRemoveWorktreeOutput_FailedCount(t *testing.T) {
	output := &usecase.RemoveWorktreeOutput{
		Results: []usecase.RemoveWorktreeResult{
			{Path: "a", Success: true},
			{Path: "b", Success: false},
			{Path: "c", Success: true},
			{Path: "d", Success: false},
		},
	}

	if output.FailedCount() != 2 {
		t.Errorf("expected 2 failures, got %d", output.FailedCount())
	}
}

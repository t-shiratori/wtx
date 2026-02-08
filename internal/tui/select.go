package tui

import (
	"errors"
	"wtx/internal/git"

	tea "github.com/charmbracelet/bubbletea"
)

// SelectWorktrees displays a TUI for selecting worktrees and returns the selected paths
func SelectWorktrees(worktrees []git.Worktree) ([]string, error) {

	if len(worktrees) == 0 {
		return nil, errors.New("no worktrees found")
	}

	m := NewModel(worktrees)

	p := tea.NewProgram(m)

	finalModel, err := p.Run()
	if err != nil {
		return nil, err
	}

	model, ok := finalModel.(Model)

	if !ok {
		return nil, errors.New("invalid TUI model")
	}

	paths := model.SelectedPaths()
	if len(paths) == 0 {
		return nil, errors.New("no worktrees selected")
	}

	return paths, nil
}

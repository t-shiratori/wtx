package tui

import (
	"wtx/internal/domain/worktree"
)

// Model manages the state of the worktree selection TUI
type Model struct {
	worktrees []worktree.Worktree // List of worktrees to display
	cursor    int                 // Current cursor position
	selected  map[int]bool        // Indices of selected worktrees
	quitting  bool                // TUI exit flag
}

// NewModel creates a new Model from the given worktree list
func NewModel(wt []worktree.Worktree) Model {
	return Model{
		worktrees: wt,
		selected:  make(map[int]bool),
	}
}

// SelectedPaths returns the paths of the selected worktrees
func (m Model) SelectedPaths() []string {
	var paths []string
	for i := range m.selected {
		if m.selected[i] {
			paths = append(paths, m.worktrees[i].Path)
		}
	}
	return paths
}

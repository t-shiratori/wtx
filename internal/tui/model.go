package tui

import (
	"wtx/internal/git"
)

// Model
// worktree 選択用 TUI の状態を管理する構造体
type Model struct {
	worktrees []git.Worktree // 表示対象の worktree 一覧
	cursor    int            // 現在のカーソル位置
	selected  map[int]bool   // 選択された worktree のインデックス
	quitting  bool           // TUI 終了フラグ
}

// NewModel
// 指定された worktree 一覧から新しい Model を生成する
func NewModel(wt []git.Worktree) Model {
	return Model{
		worktrees: wt,
		selected:  make(map[int]bool),
	}
}

// SelectedPaths
// 選択された worktree のパス一覧を返す
func (m Model) SelectedPaths() []string {
	var paths []string
	for i := range m.selected {
		if m.selected[i] {
			paths = append(paths, m.worktrees[i].Path)
		}
	}
	return paths
}

package tui

import (
	"fmt"
)

func (m Model) View() string {
	if m.quitting {
		return ""
	}

	s := "Select worktrees to remove\n\n"

	for i, wt := range m.worktrees {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if m.selected[i] {
			checked = "x"
		}

		s += fmt.Sprintf(
			"%s [%s] %s (%s)\n",
			cursor,
			checked,
			wt.Path,
			wt.Branch,
		)
	}

	s += "\n↑/↓ move • space select • enter confirm • q cancel\n"
	return s
}

package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	cursorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("14")) // bright cyan
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))  // bright red
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

		line := fmt.Sprintf(
			"%s [%s] %s (%s)",
			cursor,
			checked,
			wt.Path,
			wt.Branch,
		)

		switch {
		case m.cursor == i:
			line = cursorStyle.Render(line)
		case m.selected[i]:
			line = selectedStyle.Render(line)
		}

		s += line + "\n"
	}

	s += "\n↑/↓ move • space select • enter confirm • q cancel\n"
	return s
}

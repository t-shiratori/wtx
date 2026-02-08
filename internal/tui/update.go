package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q", "esc":
			m.quitting = true
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.worktrees)-1 {
				m.cursor++
			}

		case " ":
			// Toggle selection
			m.selected[m.cursor] = !m.selected[m.cursor]

		case "enter":
			m.quitting = true
			return m, tea.Quit
		}
	}

	return m, nil
}

package cmd

import (
	"wtx/internal/port"
	"wtx/internal/shared/spinner"
)

// hookWithSpinner wraps a HookRunner to pause/resume a spinner around hook execution.
type hookWithSpinner struct {
	inner   port.HookRunner
	spinner *spinner.Spinner
}

func newHookWithSpinner(inner port.HookRunner, sp *spinner.Spinner) *hookWithSpinner {
	return &hookWithSpinner{inner: inner, spinner: sp}
}

func (h *hookWithSpinner) Run(commands []string, dir string) error {
	if len(commands) == 0 {
		return nil
	}

	h.spinner.Pause()
	err := h.inner.Run(commands, dir)
	h.spinner.Resume()

	return err
}

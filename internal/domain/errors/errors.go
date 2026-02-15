package errors

import "fmt"

// WorktreeError represents a domain error
type WorktreeError struct {
	Op      string
	Path    string
	Message string
	Err     error
}

func (e *WorktreeError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Op, e.Path, e.Err)
	}
	return fmt.Sprintf("%s: %s: %s", e.Op, e.Path, e.Message)
}

func (e *WorktreeError) Unwrap() error {
	return e.Err
}

// NotFoundError indicates a worktree was not found
type NotFoundError struct {
	Path string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("worktree not found: %s", e.Path)
}

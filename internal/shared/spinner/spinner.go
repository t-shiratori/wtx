package spinner

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"time"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorPurple = "\033[35m"
)

var frames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// Spinner displays an animated spinner in the terminal.
type Spinner struct {
	w      io.Writer
	msg    string
	done   chan struct{}
	mu     sync.Mutex
	paused bool
}

// New creates a new Spinner that writes to w with the given message.
func New(w io.Writer, msg string) *Spinner {
	return &Spinner{
		w:    w,
		msg:  msg,
		done: make(chan struct{}),
	}
}

// Start begins the spinner animation in a goroutine.
func (s *Spinner) Start() {
	go func() {
		ticker := time.NewTicker(80 * time.Millisecond)
		defer ticker.Stop()

		i := 0
		for {
			select {
			case <-s.done:
				return
			case <-ticker.C:
				s.mu.Lock()
				if !s.paused {
					frame := frames[i%len(frames)]
					fmt.Fprintf(s.w, "\r%s%s %s%s", colorPurple, frame, s.msg, colorReset)
					i++
				}
				s.mu.Unlock()
			}
		}
	}()
}

// Pause stops the animation temporarily and clears the line.
func (s *Spinner) Pause() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.paused = true
	s.clearLine()
}

// Resume restarts the animation after a pause.
func (s *Spinner) Resume() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.paused = false
}

// clearLine clears the current line.
func (s *Spinner) clearLine() {
	fmt.Fprintf(s.w, "\r%s\r", strings.Repeat(" ", len(s.msg)+4))
}

// Stop stops the spinner and clears the line.
func (s *Spinner) Stop() {
	close(s.done)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clearLine()
}

// StopWithSuccess stops the spinner and shows a success message.
func (s *Spinner) StopWithSuccess() {
	close(s.done)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clearLine()
	fmt.Fprintf(s.w, "%s✓ %s%s\n", colorGreen, s.msg, colorReset)
}

// StopWithError stops the spinner and shows an error message.
func (s *Spinner) StopWithError() {
	close(s.done)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clearLine()
	fmt.Fprintf(s.w, "%s✗ %s%s\n", colorRed, s.msg, colorReset)
}

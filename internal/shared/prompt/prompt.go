package prompt

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Confirm writes a [y/N] prompt to w and reads the answer from r.
// It returns true if the user enters "y" or "yes" (case-insensitive).
func Confirm(w io.Writer, r io.Reader, message string) (bool, error) {
	fmt.Fprintf(w, "%s [y/N]: ", message)
	reader := bufio.NewReader(r)
	answer, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	answer = strings.TrimSpace(strings.ToLower(answer))
	return answer == "y" || answer == "yes", nil
}

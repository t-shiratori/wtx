package logger

import (
	"fmt"
	"io"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
)

func Info(w io.Writer, msg string, args ...any) {
	fmt.Fprintln(w, colorPurple+"[Info] "+fmt.Sprintf(msg, args...)+colorReset)
}

func Success(w io.Writer, msg string, args ...any) {
	fmt.Fprintln(w, colorCyan+"[Success] "+fmt.Sprintf(msg, args...)+colorReset)
}

func Warn(w io.Writer, msg string, args ...any) {
	fmt.Fprintln(w, colorYellow+"[Warn] "+fmt.Sprintf(msg, args...)+colorReset)
}

func Error(w io.Writer, msg string, args ...any) {
	fmt.Fprintln(w, colorRed+"[Error] "+fmt.Sprintf(msg, args...)+colorReset)
}

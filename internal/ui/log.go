package ui

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
)

func Info(w io.Writer, msg string, args ...any) {
	fmt.Fprintf(w, "\033[34m[Info] "+msg+"\033[0m\n", args...)
}

func Success(w io.Writer, msg string, args ...any) {
	fmt.Fprintln(w, colorGreen+"[Success] "+fmt.Sprintf(msg, args...)+colorReset)
}

func Warn(w io.Writer, msg string, args ...any) {
	fmt.Fprintln(w, colorYellow+"[Warn] "+fmt.Sprintf(msg, args...)+colorReset)
}

func Error(w io.Writer, msg string, args ...any) {
	fmt.Fprintln(w, colorRed+"[Error] "+fmt.Sprintf(msg, args...)+colorReset)
}

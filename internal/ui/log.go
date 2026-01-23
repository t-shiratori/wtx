package ui

import "fmt"

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
)

func Info(msg string, args ...any) {
	fmt.Println(colorBlue + "[Info] " + fmt.Sprintf(msg, args...) + colorReset)
}

func Success(msg string, args ...any) {
	fmt.Println(colorGreen + "[Success] " + fmt.Sprintf(msg, args...) + colorReset)
}

func Warn(msg string, args ...any) {
	fmt.Println(colorYellow + "[Warn] " + fmt.Sprintf(msg, args...) + colorReset)
}

func Error(msg string, args ...any) {
	fmt.Println(colorRed + "[Error] " + fmt.Sprintf(msg, args...) + colorReset)
}

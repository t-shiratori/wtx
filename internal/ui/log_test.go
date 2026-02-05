package ui

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

/*
自動テスト用：
- 出力にプレフィックスが含まれているか
- メッセージが正しく展開されているか
（色そのものは検証しない）
*/

func TestInfo_Format(t *testing.T) {
	var buf bytes.Buffer

	Info(&buf, "hello %s", "world")

	out := buf.String()
	if !strings.Contains(out, "[Info]") {
		t.Fatalf("expected [Info] prefix, got: %q", out)
	}
	if !strings.Contains(out, "hello world") {
		t.Fatalf("expected message, got: %q", out)
	}
}

func TestSuccess_Format(t *testing.T) {
	var buf bytes.Buffer

	Success(&buf, "ok %d", 200)

	out := buf.String()
	if !strings.Contains(out, "[Success]") {
		t.Fatalf("expected [Success] prefix, got: %q", out)
	}
}

func TestWarn_Format(t *testing.T) {
	var buf bytes.Buffer

	Warn(&buf, "warn")

	out := buf.String()
	if !strings.Contains(out, "[Warn]") {
		t.Fatalf("expected [Warn] prefix, got: %q", out)
	}
}

func TestError_Format(t *testing.T) {
	var buf bytes.Buffer

	Error(&buf, "error")

	out := buf.String()
	if !strings.Contains(out, "[Error]") {
		t.Fatalf("expected [Error] prefix, got: %q", out)
	}
}

/*
目視確認用：
- go test -v ./internal/ui
- iTerm2 / 実ターミナルで実行すること
*/

func TestLogColors_Visual(t *testing.T) {
	Info(os.Stdout, "info message")
	Success(os.Stdout, "success message")
	Warn(os.Stdout, "warn message")
	Error(os.Stdout, "error message")
}

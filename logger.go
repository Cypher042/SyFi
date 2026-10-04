package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Logger struct {
	file *os.File
}

func NewLogger() *Logger {
	if err := os.MkdirAll(appLogDirectory(), 0o755); err != nil {
		return &Logger{}
	}
	path := filepath.Join(appLogDirectory(), LogFileName)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return &Logger{}
	}
	return &Logger{file: f}
}

func (l *Logger) Close() {
	if l != nil && l.file != nil {
		_ = l.file.Close()
	}
}

func (l *Logger) log(level, msg string) {
	if l == nil || l.file == nil {
		return
	}
	clean := sanitizeForLog(msg)
	line := fmt.Sprintf("[%s] %s %s\n", strings.ToUpper(level), time.Now().Format(time.RFC3339), clean)
	_, _ = l.file.WriteString(line)
	_ = l.file.Sync()
}

func (l *Logger) Info(msg string)  { l.log("INFO", msg) }
func (l *Logger) Warn(msg string)  { l.log("WARN", msg) }
func (l *Logger) Error(msg string) { l.log("ERROR", msg) }

func sanitizeForLog(msg string) string {
	msg = strings.ReplaceAll(msg, "\n", " ")
	msg = strings.ReplaceAll(msg, "\r", " ")
	return msg
}

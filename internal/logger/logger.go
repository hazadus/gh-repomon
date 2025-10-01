// Package logger provides simple structured logging for stderr output
package logger

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Logger provides structured logging to stderr
type Logger struct {
	output  io.Writer
	verbose bool
}

// New creates a new Logger instance
func New() *Logger {
	return &Logger{
		output:  os.Stderr,
		verbose: false,
	}
}

// SetVerbose enables or disables verbose logging
func (l *Logger) SetVerbose(verbose bool) {
	l.verbose = verbose
}

// Debug logs a debug message (only shown in verbose mode)
func (l *Logger) Debug(message string) {
	if l.verbose {
		l.log("DEBUG", message, "🔧 ")
	}
}

// Info logs an informational message
func (l *Logger) Info(message string) {
	l.log("INFO", message, "")
}

// Warning logs a warning message
func (l *Logger) Warning(message string) {
	l.log("WARN", message, "⚠️ ")
}

// Error logs an error message
func (l *Logger) Error(message string) {
	l.log("ERROR", message, "❌ ")
}

// Progress logs a progress message with a search icon
func (l *Logger) Progress(message string) {
	fmt.Fprintf(l.output, "  🔍 %s\n", message)
}

// Success logs a success message with a checkmark
func (l *Logger) Success(message string) {
	fmt.Fprintf(l.output, "  ✅ %s\n", message)
}

// log is the internal logging function
func (l *Logger) log(level, message, prefix string) {
	timestamp := time.Now().UTC().Format("2006-01-02 15:04:05")
	if prefix != "" {
		fmt.Fprintf(l.output, "[%s] %s%s %s\n", timestamp, prefix, level, message)
	} else {
		fmt.Fprintf(l.output, "[%s] %s: %s\n", timestamp, level, message)
	}
}

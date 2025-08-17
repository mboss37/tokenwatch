package utils

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

// LogLevel represents the severity of a log message
type LogLevel int

const (
	// DebugLevel logs everything
	DebugLevel LogLevel = iota
	// InfoLevel logs informational messages and above
	InfoLevel
	// WarnLevel logs warnings and above
	WarnLevel
	// ErrorLevel logs errors only
	ErrorLevel
	// FatalLevel logs fatal errors only
	FatalLevel
)

// Logger provides structured logging capabilities
type Logger struct {
	mu       sync.Mutex
	level    LogLevel
	output   io.Writer
	prefix   string
	colorize bool
}

// LogEntry represents a structured log entry
type LogEntry struct {
	Level     LogLevel               `json:"level"`
	Timestamp time.Time              `json:"timestamp"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
	Caller    string                 `json:"caller,omitempty"`
}

var (
	// DefaultLogger is the global logger instance
	DefaultLogger *Logger
	once          sync.Once
)

// InitLogger initializes the default logger
func InitLogger(level LogLevel, colorize bool) {
	once.Do(func() {
		DefaultLogger = NewLogger(level, os.Stdout, "", colorize)
	})
}

// NewLogger creates a new logger instance
func NewLogger(level LogLevel, output io.Writer, prefix string, colorize bool) *Logger {
	return &Logger{
		level:    level,
		output:   output,
		prefix:   prefix,
		colorize: colorize,
	}
}

// SetLevel sets the logging level
func (l *Logger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// log writes a log entry
func (l *Logger) log(level LogLevel, msg string, fields map[string]interface{}) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	entry := LogEntry{
		Level:     level,
		Timestamp: time.Now(),
		Message:   msg,
		Fields:    fields,
	}

	// Add caller information for errors
	if level >= ErrorLevel {
		_, file, line, ok := runtime.Caller(2)
		if ok {
			parts := strings.Split(file, "/")
			if len(parts) > 1 {
				file = parts[len(parts)-1]
			}
			entry.Caller = fmt.Sprintf("%s:%d", file, line)
		}
	}

	// Format and write the log entry
	l.formatAndWrite(entry)
}

// formatAndWrite formats and writes a log entry
func (l *Logger) formatAndWrite(entry LogEntry) {
	var output string

	// Level and timestamp
	levelStr := l.levelString(entry.Level)
	timestamp := entry.Timestamp.Format("2006-01-02 15:04:05.000")

	if l.colorize {
		levelStr = l.colorizeLevel(entry.Level, levelStr)
	}

	// Base format
	output = fmt.Sprintf("[%s] %s %s", timestamp, levelStr, entry.Message)

	// Add fields if any
	if len(entry.Fields) > 0 {
		var fieldPairs []string
		for k, v := range entry.Fields {
			fieldPairs = append(fieldPairs, fmt.Sprintf("%s=%v", k, v))
		}
		output += fmt.Sprintf(" {%s}", strings.Join(fieldPairs, " "))
	}

	// Add caller for errors
	if entry.Caller != "" {
		output += fmt.Sprintf(" [%s]", entry.Caller)
	}

	// Write to output
	fmt.Fprintln(l.output, output)
}

// levelString returns the string representation of a log level
func (l *Logger) levelString(level LogLevel) string {
	switch level {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// colorizeLevel adds ANSI color codes to level strings
func (l *Logger) colorizeLevel(level LogLevel, levelStr string) string {
	switch level {
	case DebugLevel:
		return "\033[36m" + levelStr + "\033[0m" // Cyan
	case InfoLevel:
		return "\033[32m" + levelStr + "\033[0m" // Green
	case WarnLevel:
		return "\033[33m" + levelStr + "\033[0m" // Yellow
	case ErrorLevel:
		return "\033[31m" + levelStr + "\033[0m" // Red
	case FatalLevel:
		return "\033[35m" + levelStr + "\033[0m" // Magenta
	default:
		return levelStr
	}
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, fields ...map[string]interface{}) {
	l.log(DebugLevel, msg, mergeFields(fields...))
}

// Info logs an informational message
func (l *Logger) Info(msg string, fields ...map[string]interface{}) {
	l.log(InfoLevel, msg, mergeFields(fields...))
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, fields ...map[string]interface{}) {
	l.log(WarnLevel, msg, mergeFields(fields...))
}

// Error logs an error message
func (l *Logger) Error(msg string, fields ...map[string]interface{}) {
	l.log(ErrorLevel, msg, mergeFields(fields...))
}

// Fatal logs a fatal error message and exits
func (l *Logger) Fatal(msg string, fields ...map[string]interface{}) {
	l.log(FatalLevel, msg, mergeFields(fields...))
	os.Exit(1)
}

// WithFields creates a new logger with default fields
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	newLogger := &Logger{
		level:    l.level,
		output:   l.output,
		prefix:   l.prefix,
		colorize: l.colorize,
	}
	return newLogger
}

// mergeFields merges multiple field maps
func mergeFields(fields ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, f := range fields {
		for k, v := range f {
			result[k] = v
		}
	}
	return result
}

// Global logger functions

// Debug logs a debug message using the default logger
func Debug(msg string, fields ...map[string]interface{}) {
	if DefaultLogger != nil {
		DefaultLogger.Debug(msg, fields...)
	}
}

// Info logs an informational message using the default logger
func Info(msg string, fields ...map[string]interface{}) {
	if DefaultLogger != nil {
		DefaultLogger.Info(msg, fields...)
	}
}

// Warn logs a warning message using the default logger
func Warn(msg string, fields ...map[string]interface{}) {
	if DefaultLogger != nil {
		DefaultLogger.Warn(msg, fields...)
	}
}

// Error logs an error message using the default logger
func Error(msg string, fields ...map[string]interface{}) {
	if DefaultLogger != nil {
		DefaultLogger.Error(msg, fields...)
	}
}

// Fatal logs a fatal error message using the default logger and exits
func Fatal(msg string, fields ...map[string]interface{}) {
	if DefaultLogger != nil {
		DefaultLogger.Fatal(msg, fields...)
	} else {
		fmt.Fprintf(os.Stderr, "FATAL: %s\n", msg)
		os.Exit(1)
	}
}

// ParseLogLevel parses a string log level
func ParseLogLevel(level string) LogLevel {
	switch strings.ToLower(level) {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn", "warning":
		return WarnLevel
	case "error":
		return ErrorLevel
	case "fatal":
		return FatalLevel
	default:
		return InfoLevel
	}
}

// Package logger provides structured logging functionality.
package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
)

// Level represents the log level
type Level slog.Level

// Log levels
const (
	LevelDebug Level = Level(slog.LevelDebug)
	LevelInfo  Level = Level(slog.LevelInfo)
	LevelWarn  Level = Level(slog.LevelWarn)
	LevelError Level = Level(slog.LevelError)
)

// Logger wraps slog.Logger with additional functionality
type Logger struct {
	logger *slog.Logger
}

// Config holds logger configuration
type Config struct {
	Level  Level
	Format Format
	Output io.Writer
}

// Format represents the log output format
type Format string

// Format types
const (
	FormatJSON Format = "json"
	FormatText Format = "text"
)

// New creates a new logger with the given configuration
func New(config Config) *Logger {
	if config.Output == nil {
		config.Output = os.Stderr
	}

	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level: slog.Level(config.Level),
	}

	switch config.Format {
	case FormatJSON:
		handler = slog.NewJSONHandler(config.Output, opts)
	case FormatText:
		handler = slog.NewTextHandler(config.Output, opts)
	default:
		handler = slog.NewTextHandler(config.Output, opts)
	}

	return &Logger{
		logger: slog.New(handler),
	}
}

// Default creates a logger with default configuration
func Default() *Logger {
	return New(Config{
		Level:  LevelInfo,
		Format: FormatText,
		Output: os.Stderr,
	})
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

// Info logs an info message
func (l *Logger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

// Error logs an error message
func (l *Logger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

// DebugContext logs a debug message with context
func (l *Logger) DebugContext(ctx context.Context, msg string, args ...any) {
	l.logger.DebugContext(ctx, msg, args...)
}

// InfoContext logs an info message with context
func (l *Logger) InfoContext(ctx context.Context, msg string, args ...any) {
	l.logger.InfoContext(ctx, msg, args...)
}

// WarnContext logs a warning message with context
func (l *Logger) WarnContext(ctx context.Context, msg string, args ...any) {
	l.logger.WarnContext(ctx, msg, args...)
}

// ErrorContext logs an error message with context
func (l *Logger) ErrorContext(ctx context.Context, msg string, args ...any) {
	l.logger.ErrorContext(ctx, msg, args...)
}

// With returns a new logger with the given attributes
func (l *Logger) With(args ...any) *Logger {
	return &Logger{
		logger: l.logger.With(args...),
	}
}

// WithGroup returns a new logger with the given group name
func (l *Logger) WithGroup(name string) *Logger {
	return &Logger{
		logger: l.logger.WithGroup(name),
	}
}

// Handler returns the underlying slog handler
func (l *Logger) Handler() slog.Handler {
	return l.logger.Handler()
}

// Enabled reports whether the logger emits a log record at the given level
func (l *Logger) Enabled(ctx context.Context, level Level) bool {
	return l.logger.Enabled(ctx, slog.Level(level))
}

// Global logger instance
var global = Default()

// SetGlobal sets the global logger
func SetGlobal(logger *Logger) {
	global = logger
}

// Global returns the global logger
func Global() *Logger {
	return global
}

// Debug logs a debug message using the global logger
func Debug(msg string, args ...any) {
	global.Debug(msg, args...)
}

// Info logs an info message using the global logger
func Info(msg string, args ...any) {
	global.Info(msg, args...)
}

// Warn logs a warning message using the global logger
func Warn(msg string, args ...any) {
	global.Warn(msg, args...)
}

// Error logs an error message using the global logger
func Error(msg string, args ...any) {
	global.Error(msg, args...)
}

// DebugContext logs a debug message with context using the global logger
func DebugContext(ctx context.Context, msg string, args ...any) {
	global.DebugContext(ctx, msg, args...)
}

// InfoContext logs an info message with context using the global logger
func InfoContext(ctx context.Context, msg string, args ...any) {
	global.InfoContext(ctx, msg, args...)
}

// WarnContext logs a warning message with context using the global logger
func WarnContext(ctx context.Context, msg string, args ...any) {
	global.WarnContext(ctx, msg, args...)
}

// ErrorContext logs an error message with context using the global logger
func ErrorContext(ctx context.Context, msg string, args ...any) {
	global.ErrorContext(ctx, msg, args...)
}

// With returns a new logger with the given attributes using the global logger
func With(args ...any) *Logger {
	return global.With(args...)
}

// WithGroup returns a new logger with the given group name using the global logger
func WithGroup(name string) *Logger {
	return global.WithGroup(name)
}
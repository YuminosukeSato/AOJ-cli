package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("JSON format", func(t *testing.T) {
		buf := &bytes.Buffer{}
		logger := New(Config{
			Level:  LevelDebug,
			Format: FormatJSON,
			Output: buf,
		})

		logger.Info("test message", "key", "value")

		output := buf.String()
		assert.Contains(t, output, "test message")
		assert.Contains(t, output, "key")
		assert.Contains(t, output, "value")

		// Verify it's valid JSON
		var jsonData map[string]interface{}
		err := json.Unmarshal([]byte(output), &jsonData)
		assert.NoError(t, err)
	})

	t.Run("Text format", func(t *testing.T) {
		buf := &bytes.Buffer{}
		logger := New(Config{
			Level:  LevelDebug,
			Format: FormatText,
			Output: buf,
		})

		logger.Info("test message", "key", "value")

		output := buf.String()
		assert.Contains(t, output, "test message")
		assert.Contains(t, output, "key=value")
	})

	t.Run("Default format when unspecified", func(t *testing.T) {
		buf := &bytes.Buffer{}
		logger := New(Config{
			Level:  LevelDebug,
			Format: "",
			Output: buf,
		})

		logger.Info("test message")

		output := buf.String()
		assert.Contains(t, output, "test message")
		// Should use text format by default
		assert.Contains(t, output, "level=INFO")
	})
}

func TestDefault(t *testing.T) {
	logger := Default()
	assert.NotNil(t, logger)
	assert.NotNil(t, logger.logger)
}

func TestLoggingMethods(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Config{
		Level:  LevelDebug,
		Format: FormatText,
		Output: buf,
	})

	t.Run("Debug", func(t *testing.T) {
		buf.Reset()
		logger.Debug("debug message", "key", "value")
		output := buf.String()
		assert.Contains(t, output, "debug message")
		assert.Contains(t, output, "level=DEBUG")
	})

	t.Run("Info", func(t *testing.T) {
		buf.Reset()
		logger.Info("info message", "key", "value")
		output := buf.String()
		assert.Contains(t, output, "info message")
		assert.Contains(t, output, "level=INFO")
	})

	t.Run("Warn", func(t *testing.T) {
		buf.Reset()
		logger.Warn("warn message", "key", "value")
		output := buf.String()
		assert.Contains(t, output, "warn message")
		assert.Contains(t, output, "level=WARN")
	})

	t.Run("Error", func(t *testing.T) {
		buf.Reset()
		logger.Error("error message", "key", "value")
		output := buf.String()
		assert.Contains(t, output, "error message")
		assert.Contains(t, output, "level=ERROR")
	})
}

func TestContextMethods(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Config{
		Level:  LevelDebug,
		Format: FormatText,
		Output: buf,
	})

	ctx := context.Background()

	t.Run("DebugContext", func(t *testing.T) {
		buf.Reset()
		logger.DebugContext(ctx, "debug context message")
		output := buf.String()
		assert.Contains(t, output, "debug context message")
		assert.Contains(t, output, "level=DEBUG")
	})

	t.Run("InfoContext", func(t *testing.T) {
		buf.Reset()
		logger.InfoContext(ctx, "info context message")
		output := buf.String()
		assert.Contains(t, output, "info context message")
		assert.Contains(t, output, "level=INFO")
	})

	t.Run("WarnContext", func(t *testing.T) {
		buf.Reset()
		logger.WarnContext(ctx, "warn context message")
		output := buf.String()
		assert.Contains(t, output, "warn context message")
		assert.Contains(t, output, "level=WARN")
	})

	t.Run("ErrorContext", func(t *testing.T) {
		buf.Reset()
		logger.ErrorContext(ctx, "error context message")
		output := buf.String()
		assert.Contains(t, output, "error context message")
		assert.Contains(t, output, "level=ERROR")
	})
}

func TestWith(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Config{
		Level:  LevelDebug,
		Format: FormatText,
		Output: buf,
	})

	childLogger := logger.With("component", "test", "version", "1.0")
	childLogger.Info("test message")

	output := buf.String()
	assert.Contains(t, output, "test message")
	assert.Contains(t, output, "component=test")
	assert.Contains(t, output, "version=1.0")
}

func TestWithGroup(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Config{
		Level:  LevelDebug,
		Format: FormatJSON,
		Output: buf,
	})

	groupLogger := logger.WithGroup("database")
	groupLogger.Info("connection established", "host", "localhost")

	output := buf.String()
	assert.Contains(t, output, "connection established")
	
	// Parse JSON to verify group structure
	var jsonData map[string]interface{}
	err := json.Unmarshal([]byte(output), &jsonData)
	assert.NoError(t, err)
	
	// Check if the database group exists
	database, exists := jsonData["database"]
	assert.True(t, exists)
	assert.IsType(t, map[string]interface{}{}, database)
}

func TestEnabled(t *testing.T) {
	logger := New(Config{
		Level:  LevelWarn,
		Format: FormatText,
		Output: &bytes.Buffer{},
	})

	ctx := context.Background()

	assert.False(t, logger.Enabled(ctx, LevelDebug))
	assert.False(t, logger.Enabled(ctx, LevelInfo))
	assert.True(t, logger.Enabled(ctx, LevelWarn))
	assert.True(t, logger.Enabled(ctx, LevelError))
}

func TestLogLevel(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Config{
		Level:  LevelWarn,
		Format: FormatText,
		Output: buf,
	})

	// These should not be logged due to log level
	logger.Debug("debug message")
	logger.Info("info message")
	
	// These should be logged
	logger.Warn("warn message")
	logger.Error("error message")

	output := buf.String()
	assert.NotContains(t, output, "debug message")
	assert.NotContains(t, output, "info message")
	assert.Contains(t, output, "warn message")
	assert.Contains(t, output, "error message")
}

func TestGlobalLogger(t *testing.T) {
	// Save original global logger
	originalGlobal := Global()
	defer SetGlobal(originalGlobal)

	buf := &bytes.Buffer{}
	logger := New(Config{
		Level:  LevelDebug,
		Format: FormatText,
		Output: buf,
	})

	SetGlobal(logger)
	assert.Equal(t, logger, Global())

	// Test global functions
	Info("global info message")
	output := buf.String()
	assert.Contains(t, output, "global info message")

	buf.Reset()
	Error("global error message", "key", "value")
	output = buf.String()
	assert.Contains(t, output, "global error message")
	assert.Contains(t, output, "key=value")
}

func TestGlobalFunctions(t *testing.T) {
	// Save original global logger
	originalGlobal := Global()
	defer SetGlobal(originalGlobal)

	buf := &bytes.Buffer{}
	logger := New(Config{
		Level:  LevelDebug,
		Format: FormatText,
		Output: buf,
	})
	SetGlobal(logger)

	ctx := context.Background()

	t.Run("Global debug functions", func(t *testing.T) {
		buf.Reset()
		Debug("global debug")
		DebugContext(ctx, "global debug context")
		output := buf.String()
		lines := strings.Split(strings.TrimSpace(output), "\n")
		assert.Len(t, lines, 2)
		assert.Contains(t, output, "global debug")
		assert.Contains(t, output, "global debug context")
	})

	t.Run("Global With functions", func(t *testing.T) {
		childLogger := With("global", "true")
		assert.NotNil(t, childLogger)

		groupLogger := WithGroup("global")
		assert.NotNil(t, groupLogger)
	})
}
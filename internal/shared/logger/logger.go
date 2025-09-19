package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Custom log levels
const (
	TraceLevel   zapcore.Level = -3 // Very detailed debugging
	VerboseLevel zapcore.Level = -2 // More detailed than debug
)

var (
	// Logger is the global logger instance
	Logger *zap.Logger

	// Sugar provides a more ergonomic API
	Sugar *zap.SugaredLogger

	// logLevels defines available log levels
	logLevels = map[string]zapcore.Level{
		"trace":   TraceLevel,
		"verbose": VerboseLevel,
		"debug":   zapcore.DebugLevel,
		"info":    zapcore.InfoLevel,
		"warn":    zapcore.WarnLevel,
		"error":   zapcore.ErrorLevel,
		"fatal":   zapcore.FatalLevel,
	}

	// Custom level names for display
	levelNames = map[zapcore.Level]string{
		TraceLevel:   "TRACE",
		VerboseLevel: "VERBOSE",
	}

	once sync.Once
	mu   sync.RWMutex
)

// LogConfig holds logger configuration
type LogConfig struct {
	Level      string
	Format     string // "console", "json"
	Output     string // "stderr", "stdout", or file path
	NoColor    bool
	ShowCaller bool
	ShowTime   bool
}

// Initialize sets up the logger with the given configuration
func Initialize(cfg *LogConfig) error {
	var err error
	once.Do(func() {
		err = initializeLogger(cfg)
	})
	return err
}

func initializeLogger(cfg *LogConfig) error {
	mu.Lock()
	defer mu.Unlock()

	if cfg == nil {
		cfg = &LogConfig{
			Level:      getLogLevel(),
			Format:     "console",
			Output:     "stderr",
			NoColor:    false,
			ShowCaller: false,
			ShowTime:   true,
		}
	}

	// Parse log level
	level, exists := logLevels[strings.ToLower(cfg.Level)]
	if !exists {
		level = zapcore.InfoLevel
	}

	// Create encoder config
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    customLevelEncoder(cfg.NoColor),
		EncodeTime:     customTimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   customCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	// Create encoder based on format
	var encoder zapcore.Encoder
	switch cfg.Format {
	case "json":
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	default:
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// Create output writer
	var output zapcore.WriteSyncer
	switch cfg.Output {
	case "stdout":
		output = zapcore.AddSync(os.Stdout)
	case "stderr":
		output = zapcore.AddSync(os.Stderr)
	default:
		// File output
		file, err := os.OpenFile(cfg.Output, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		output = zapcore.AddSync(file)
	}

	// Create core
	core := zapcore.NewCore(encoder, output, level)

	// Build logger with options
	opts := []zap.Option{
		zap.AddCallerSkip(1), // Skip logger wrapper functions
	}

	if cfg.ShowCaller {
		opts = append(opts, zap.AddCaller())
	}

	// Add custom level support
	core = zapcore.NewTee(
		core,
		&customCore{Core: core},
	)

	Logger = zap.New(core, opts...)
	Sugar = Logger.Sugar()

	// Log initialization - use the new logger directly
	Logger.Info("Logger initialized",
		zap.String("level", cfg.Level),
		zap.String("format", cfg.Format),
		zap.String("output", cfg.Output),
		zap.Bool("color", !cfg.NoColor),
	)

	return nil
}

// customCore handles custom log levels
type customCore struct {
	zapcore.Core
}

func (c *customCore) Check(entry zapcore.Entry, checked *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return checked.AddCore(entry, c)
	}
	return checked
}

// Custom encoders
func customLevelEncoder(noColor bool) zapcore.LevelEncoder {
	return func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		var levelStr string

		// Check custom levels first
		if name, ok := levelNames[level]; ok {
			levelStr = name
		} else {
			levelStr = level.CapitalString()
		}

		if noColor {
			enc.AppendString(levelStr)
			return
		}

		// Color based on level
		var coloredLevel string
		switch level {
		case TraceLevel:
			coloredLevel = color.HiBlackString(levelStr)
		case VerboseLevel:
			coloredLevel = color.HiBlueString(levelStr)
		case zapcore.DebugLevel:
			coloredLevel = color.CyanString(levelStr)
		case zapcore.InfoLevel:
			coloredLevel = color.GreenString(levelStr)
		case zapcore.WarnLevel:
			coloredLevel = color.YellowString(levelStr)
		case zapcore.ErrorLevel:
			coloredLevel = color.RedString(levelStr)
		case zapcore.FatalLevel:
			coloredLevel = color.HiRedString(levelStr)
		default:
			coloredLevel = levelStr
		}

		enc.AppendString(coloredLevel)
	}
}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("15:04:05.000"))
}

func customCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	// Show only filename and line number
	enc.AppendString(fmt.Sprintf("%s:%d",
		filepath.Base(caller.File),
		caller.Line))
}

// Helper functions for structured logging

// getZapFields converts key-value pairs to zap fields
func getZapFields(keyvals ...interface{}) []zap.Field {
	fields := make([]zap.Field, 0, len(keyvals)/2)
	for i := 0; i < len(keyvals)-1; i += 2 {
		key, ok := keyvals[i].(string)
		if !ok {
			continue
		}
		fields = append(fields, zap.Any(key, keyvals[i+1]))
	}
	return fields
}

// Logging functions

// Trace logs trace-level messages (very detailed debugging)
func Trace(msg string, keyvals ...interface{}) {
	if Logger == nil {
		fmt.Fprintf(os.Stderr, "[TRACE] %s\n", msg)
		return
	}
	if Logger.Core().Enabled(TraceLevel) {
		fields := getZapFields(keyvals...)
		fields = append(fields, zap.String("_level", "TRACE"))
		Logger.Debug(msg, fields...)
	}
}

// Verbose logs verbose-level messages (more than debug, less than trace)
func Verbose(msg string, keyvals ...interface{}) {
	if Logger == nil {
		fmt.Fprintf(os.Stderr, "[VERBOSE] %s\n", msg)
		return
	}
	if Logger.Core().Enabled(VerboseLevel) {
		fields := getZapFields(keyvals...)
		fields = append(fields, zap.String("_level", "VERBOSE"))
		Logger.Debug(msg, fields...)
	}
}

// Debug logs debug messages
func Debug(msg string, keyvals ...interface{}) {
	if Logger == nil {
		fmt.Fprintf(os.Stderr, "[DEBUG] %s\n", msg)
		return
	}
	Logger.Debug(msg, getZapFields(keyvals...)...)
}

// Info logs informational messages
func Info(msg string, keyvals ...interface{}) {
	if Logger == nil {
		fmt.Fprintf(os.Stderr, "[INFO] %s\n", msg)
		return
	}
	Logger.Info(msg, getZapFields(keyvals...)...)
}

// Warn logs warning messages
func Warn(msg string, keyvals ...interface{}) {
	if Logger == nil {
		fmt.Fprintf(os.Stderr, "[WARN] %s\n", msg)
		return
	}
	Logger.Warn(msg, getZapFields(keyvals...)...)
}

// Error logs error messages
func Error(msg string, keyvals ...interface{}) {
	if Logger == nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", msg)
		return
	}
	Logger.Error(msg, getZapFields(keyvals...)...)
}

// Fatal logs fatal messages and exits
func Fatal(msg string, keyvals ...interface{}) {
	if Logger == nil {
		fmt.Fprintf(os.Stderr, "[FATAL] %s\n", msg)
		os.Exit(1)
		return
	}
	Logger.Fatal(msg, getZapFields(keyvals...)...)
}

// WithFields returns a logger with persistent fields
func WithFields(keyvals ...interface{}) *zap.Logger {
	if Logger == nil {
		return nil
	}
	return Logger.With(getZapFields(keyvals...)...)
}

// Command-specific logging helpers

// LogCommand logs command execution details
func LogCommand(cmd string, args []string, keyvals ...interface{}) {
	// Convert fields back to key-value pairs for Debug function
	allKeyvals := append(keyvals,
		"command", cmd,
		"args", args,
	)
	Debug("Executing command", allKeyvals...)
}

// LogDuration logs operation duration
func LogDuration(operation string, start time.Time, keyvals ...interface{}) {
	duration := time.Since(start)
	fields := getZapFields(keyvals...)
	fields = append(fields,
		zap.Duration("duration", duration),
		zap.String("operation", operation),
	)

	if duration > time.Second {
		if Logger == nil {
			fmt.Fprintf(os.Stderr, "[INFO] Operation completed: %s (duration: %v)\n", operation, duration)
			return
		}
		Logger.Info("Operation completed", fields...)
	} else {
		if Logger == nil {
			fmt.Fprintf(os.Stderr, "[DEBUG] Operation completed: %s (duration: %v)\n", operation, duration)
			return
		}
		Logger.Debug("Operation completed", fields...)
	}
}

// LogError logs an error with context
func LogError(err error, msg string, keyvals ...interface{}) {
	if err == nil {
		return
	}

	fields := getZapFields(keyvals...)
	fields = append(fields, zap.Error(err))

	if Logger == nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s: %v\n", msg, err)
		return
	}
	Logger.Error(msg, fields...)
}

// Helper to get log level from environment
func getLogLevel() string {
	// Priority: CLI_T_LOG_LEVEL > LOG_LEVEL > default
	if level := os.Getenv("CLI_T_LOG_LEVEL"); level != "" {
		return level
	}
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		return level
	}

	// Check debug mode
	if os.Getenv("CLI_T_DEBUG") == "true" {
		return "debug"
	}

	return "info"
}

// SetLevel dynamically changes the log level
func SetLevel(level string) error {
	mu.Lock()
	defer mu.Unlock()

	_, exists := logLevels[strings.ToLower(level)]
	if !exists {
		return fmt.Errorf("invalid log level: %s", level)
	}

	// This requires reinitializing the logger
	// For now, we'll just return an error since we can't reinitialize with sync.Once
	return fmt.Errorf("dynamic level change not supported after initialization")
}

// GetLevel returns the current log level
func GetLevel() string {
	mu.RLock()
	defer mu.RUnlock()

	if Logger == nil {
		return "info"
	}

	// Check which level is enabled (from highest to lowest)
	levelOrder := []struct {
		name  string
		level zapcore.Level
	}{
		{"trace", TraceLevel},
		{"verbose", VerboseLevel},
		{"debug", zapcore.DebugLevel},
		{"info", zapcore.InfoLevel},
		{"warn", zapcore.WarnLevel},
		{"error", zapcore.ErrorLevel},
		{"fatal", zapcore.FatalLevel},
	}

	for _, l := range levelOrder {
		if Logger.Core().Enabled(l.level) {
			return l.name
		}
	}

	return "info"
}

// Sync flushes any buffered log entries
func Sync() error {
	if Logger != nil {
		return Logger.Sync()
	}
	return nil
}

package storage

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"gorm.io/gorm/logger"
)

// SlogGormLogger implements gorm/logger.Interface using slog
type SlogGormLogger struct {
	logger *slog.Logger
	level  slog.Level
}

// Info logs info-level messages
func (l *SlogGormLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	if l.level <= slog.LevelInfo {
		l.logger.InfoContext(ctx, msg, "args", args)
	}
}

// Warn logs warn-level messages
func (l *SlogGormLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	if l.level <= slog.LevelWarn {
		l.logger.WarnContext(ctx, msg, "args", args)
	}
}

// Error logs error-level messages
func (l *SlogGormLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	if l.level <= slog.LevelError {
		l.logger.ErrorContext(ctx, msg, "args", args)
	}
}

// Trace logs SQL queries with execution time
func (l *SlogGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.level > slog.LevelDebug {
		return
	}
	sql, rows := fc()
	duration := time.Since(begin)

	l.logger.DebugContext(ctx, "SQL Query",
		"sql", sql,
		"rows", rows,
		"duration", duration,
		"error", err,
	)
}

// LogMode allows changing the logging level dynamically
func (l *SlogGormLogger) LogMode(level logger.LogLevel) logger.Interface {
	var newLevel slog.Level

	switch level {
	case logger.Silent:
		newLevel = slog.LevelError + 1 // Silent, suppresses all logs
	case logger.Error:
		newLevel = slog.LevelError
	case logger.Warn:
		newLevel = slog.LevelWarn
	case logger.Info:
		newLevel = slog.LevelInfo
	default:
		newLevel = slog.LevelInfo
	}

	newLogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: newLevel}))
	return &SlogGormLogger{logger: newLogger, level: newLevel}
}

func parseLogLevel(levelStr string) slog.Level {
	switch strings.ToLower(levelStr) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		fmt.Printf("Unknown log level: %s, defaulting to INFO\n", levelStr)
		return slog.LevelInfo
	}
}

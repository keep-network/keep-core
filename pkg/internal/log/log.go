// Package log implements a basic logging utility. It is currently tied to the
// zap logging implementation (see more here: https://godoc.org/go.uber.org/zap
// and here:https://github.com/uber-go/zap/blob/master/FAQ.md).
package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is an abstraction over the zap logging library, which is a configurable,
// structured logger, safe for concurrent use.
type Logger struct {
	zap *zap.Logger
}

var (
	// Any encapsulates the zap catch-all type.
	Any = zap.Any
	// Error encapsulates the zap error type.
	Error = zap.Error
)

// NewLogger returns a NewLogger which is primed for use in Development.
// TODO: make this configurable
func NewLogger() *Logger {
	logger, err := zap.NewDevelopment(
		zap.AddCaller(),
	)
	if err != nil {
		panic(err)
	}
	return &Logger{
		zap: logger,
	}
}

// Field is an alias for Field. Per, zap's documentation, aliasing this type
// dramatically improves the navigability of this package's API documentation.
type Field = zapcore.Field

// Debug is a wrapper around zap's debug.
func (l *Logger) Debug(message string, fields ...Field) {
	l.zap.Debug(message, fields...)
}

// Info is a wrapper around zap's Info.
func (l *Logger) Info(message string, fields ...Field) {
	l.zap.Info(message, fields...)
}

// Warn is a wrapper around zap's Warn.
func (l *Logger) Warn(message string, fields ...Field) {
	l.zap.Warn(message, fields...)
}

// Error is a wrapper around zap's Error.
func (l *Logger) Error(message string, fields ...Field) {
	l.zap.Error(message, fields...)
}

// Fatal is a wrapper around zap's Fatal.
func (l *Logger) Fatal(message string, fields ...Field) {
	l.zap.Fatal(message, fields...)
}

// Sync calls the underlying Core's Sync method, flushing any buffered log
// entries. Applications should take care to call Sync before exiting.
func (l *Logger) Sync() error {
	return l.zap.Sync()
}

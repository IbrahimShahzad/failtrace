// Author: Ibrahim Shahzad,
// Date: 2025-06-12,
// License: MIT
//
// Package logger provides a simple request-scoped logging mechanism that buffers log entries
// and writes them out only when an error occurs or when explicitly flushed.
//
// Each request gets a unique ID, and logs are written to stderr by default.
//
// > Note: This is *NOT* thread-safe. It is designed for request-scoped logging
//
// Usage:
//
//		func handleRequest(ctx context.Context) error {
//		    ctx = logger.WithLogger(ctx)
//		    log := logger.FromContext(ctx)
//		    defer log.FlushIf(nil)
//
//		    log.Debug("Processing request")
//		    log.Info("Validating input")
//
//		    if err := someOperation(); err != nil {
//		        log.Error("Operation failed")
//		        log.FlushIf(err)
//		        return err
//		    }
//
//	     	log.FlushIf(nil)
//		    return nil
//		}
//
// For more usage examples, see examples in the package documentation.
package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/google/uuid"
)

type ctxKey struct{}

type Level byte

const (
	DebugLevel Level = 'D'
	InfoLevel  Level = 'I'
	WarnLevel  Level = 'W'
	ErrorLevel Level = 'E'
)

type logEntry struct {
	level   Level
	message string
}

type requestLogger struct {
	id  string
	buf []logEntry
	w   io.Writer
}

var pool = sync.Pool{
	New: func() any {
		return &requestLogger{
			buf: make([]logEntry, 0, 32),
			w:   os.Stderr,
		}
	},
}

// WithLogger returns a new context with logger.
func WithLogger(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKey{}, pool.Get().(*requestLogger).reset())
}

// FromContext retrieves the logger from the context.
func FromContext(ctx context.Context) *requestLogger {
	if rl, ok := ctx.Value(ctxKey{}).(*requestLogger); ok {
		return rl
	}
	return &requestLogger{
		id:  "noop",
		buf: make([]logEntry, 0),
		w:   io.Discard,
	}
}

// Debug logs an debug-level message. takes string as input.
//
// Usage example:
//
//	logger := &requestLogger{}
//	logger.Debug("failed to process request")
func (l *requestLogger) Debug(msg string) {
	l.buf = append(l.buf, logEntry{DebugLevel, msg})
}

// Debugf logs an debug-level message.
// The message is constructed by formatting the provided arguments using fmt.Sprint.
//
// Usage example:
//
//	logger := &requestLogger{}
//	logger.Debugf("failed to process request: %v", err)
func (l *requestLogger) Debugf(format string, args ...any) {
	l.buf = append(l.buf, logEntry{DebugLevel, fmt.Sprintf(format, args...)})
}

// Info logs an info-level message. takes string as input.
//
// Usage example:
//
//	logger := &requestLogger{}
//	logger.Info("failed to process request")
func (l *requestLogger) Info(msg string) {
	l.buf = append(l.buf, logEntry{InfoLevel, msg})
}

// Infof logs an info-level message.
// The message is constructed by formatting the provided arguments using fmt.Sprint.
//
// Usage example:
//
//	logger := &requestLogger{}
//	logger.Infof("failed to process request: %v", err)
func (l *requestLogger) Infof(format string, args ...any) {
	l.buf = append(l.buf, logEntry{InfoLevel, fmt.Sprintf(format, args...)})
}

// Warn logs an warn-level message. takes string as input.
//
// Usage example:
//
//	logger := &requestLogger{}
//	logger.Warn("failed to process request")
func (l *requestLogger) Warn(msg string) {
	l.buf = append(l.buf, logEntry{WarnLevel, msg})
}

// Warnf logs an warn-level message.
// The message is constructed by formatting the provided arguments using fmt.Sprint.
//
// Usage example:
//
//	logger := &requestLogger{}
//	logger.Warnf("failed to process request: %v", err)
func (l *requestLogger) Warnf(format string, args ...any) {
	l.buf = append(l.buf, logEntry{WarnLevel, fmt.Sprintf(format, args...)})
}

// Errorf logs an error-level message.
// The message is constructed by formatting the provided arguments using fmt.Sprint.
//
// Usage example:
//
//	logger := &requestLogger{}
//	logger.Errorf("failed to process request: %v", err)
func (l *requestLogger) Errorf(format string, args ...any) {
	l.buf = append(l.buf, logEntry{ErrorLevel, fmt.Sprintf(format, args...)})
}

// Error logs an error-level message. takes string as input.
//
// Usage example:
//
//	logger := &requestLogger{}
//	logger.Error("failed to process request")
func (l *requestLogger) Error(msg string) {
	l.buf = append(l.buf, logEntry{ErrorLevel, msg}) // Should be ErrorLevel, not WarnLevel
}

// FlushIf writes buffered log entries and the given error to the writer if err is not nil,
// then returns the logger to the pool.
func (l *requestLogger) FlushIf(err error) {
	defer l.put()

	if err == nil {
		return
	}

	for _, entry := range l.buf {
		if _, wErr := fmt.Fprintf(l.w, "[%s] %c: %s\n", l.id, entry.level, entry.message); wErr != nil {
			_ = wErr
		}
	}

	if _, wErr := fmt.Fprintf(l.w, "[%s] E: %v\n", l.id, err); wErr != nil {
		_ = wErr
	}
}

// Flush writes buffered log entries, then returns the logger to the pool.
func (l *requestLogger) Flush() {
	defer l.put()

	for _, entry := range l.buf {
		if _, wErr := fmt.Fprintf(l.w, "[%s] %c: %s\n", l.id, entry.level, entry.message); wErr != nil {
			_ = wErr
		}
	}
}

// put resets the logger's buffer and ID, effectively clearing all logs.
func (l *requestLogger) put() {
	pool.Put(l.reset())
}

func (l *requestLogger) reset() *requestLogger {
	l.buf = l.buf[:0]
	l.id = uuid.New().String()
	return l
}

package logger

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"testing"
)

func TestRequestLogger_Debug(t *testing.T) {
	var buf bytes.Buffer
	logger := &requestLogger{
		id:  "test-123",
		buf: make([]logEntry, 0),
		w:   &buf,
	}

	logger.Debug("test debug message")

	if len(logger.buf) != 1 {
		t.Errorf("Expected 1 log entry, got %d", len(logger.buf))
	}
	if logger.buf[0].level != DebugLevel {
		t.Errorf("Expected DebugLevel, got %c", logger.buf[0].level)
	}
	if logger.buf[0].message != "test debug message" {
		t.Errorf("Expected 'test debug message', got '%s'", logger.buf[0].message)
	}
}

func TestRequestLogger_Debugf(t *testing.T) {
	var buf bytes.Buffer
	logger := &requestLogger{
		id:  "test-123",
		buf: make([]logEntry, 0),
		w:   &buf,
	}

	logger.Debugf("test debug %s %d", "message", 42)

	if len(logger.buf) != 1 {
		t.Errorf("Expected 1 log entry, got %d", len(logger.buf))
	}
	if logger.buf[0].level != DebugLevel {
		t.Errorf("Expected DebugLevel, got %c", logger.buf[0].level)
	}
	if logger.buf[0].message != "test debug message 42" {
		t.Errorf("Expected 'test debug message 42', got '%s'", logger.buf[0].message)
	}
}

func TestRequestLogger_Info(t *testing.T) {
	var buf bytes.Buffer
	logger := &requestLogger{
		id:  "test-123",
		buf: make([]logEntry, 0),
		w:   &buf,
	}

	logger.Info("test info message")

	if len(logger.buf) != 1 {
		t.Errorf("Expected 1 log entry, got %d", len(logger.buf))
	}
	if logger.buf[0].level != InfoLevel {
		t.Errorf("Expected InfoLevel, got %c", logger.buf[0].level)
	}
	if logger.buf[0].message != "test info message" {
		t.Errorf("Expected 'test info message', got '%s'", logger.buf[0].message)
	}
}

func TestRequestLogger_Infof(t *testing.T) {
	var buf bytes.Buffer
	logger := &requestLogger{
		id:  "test-123",
		buf: make([]logEntry, 0),
		w:   &buf,
	}

	logger.Infof("test info %s", "formatted")

	if len(logger.buf) != 1 {
		t.Errorf("Expected 1 log entry, got %d", len(logger.buf))
	}
	if logger.buf[0].level != InfoLevel {
		t.Errorf("Expected InfoLevel, got %c", logger.buf[0].level)
	}
	if logger.buf[0].message != "test info formatted" {
		t.Errorf("Expected 'test info formatted', got '%s'", logger.buf[0].message)
	}
}

func TestRequestLogger_Warn(t *testing.T) {
	var buf bytes.Buffer
	logger := &requestLogger{
		id:  "test-123",
		buf: make([]logEntry, 0),
		w:   &buf,
	}

	logger.Warn("test warn message")

	if len(logger.buf) != 1 {
		t.Errorf("Expected 1 log entry, got %d", len(logger.buf))
	}
	if logger.buf[0].level != WarnLevel {
		t.Errorf("Expected WarnLevel, got %c", logger.buf[0].level)
	}
	if logger.buf[0].message != "test warn message" {
		t.Errorf("Expected 'test warn message', got '%s'", logger.buf[0].message)
	}
}

func TestRequestLogger_Warnf(t *testing.T) {
	var buf bytes.Buffer
	logger := &requestLogger{
		id:  "test-123",
		buf: make([]logEntry, 0),
		w:   &buf,
	}

	logger.Warnf("test warn %d", 123)

	if len(logger.buf) != 1 {
		t.Errorf("Expected 1 log entry, got %d", len(logger.buf))
	}
	if logger.buf[0].level != WarnLevel {
		t.Errorf("Expected WarnLevel, got %c", logger.buf[0].level)
	}
	if logger.buf[0].message != "test warn 123" {
		t.Errorf("Expected 'test warn 123', got '%s'", logger.buf[0].message)
	}
}

func TestRequestLogger_Error(t *testing.T) {
	var buf bytes.Buffer
	logger := &requestLogger{
		id:  "test-123",
		buf: make([]logEntry, 0),
		w:   &buf,
	}

	logger.Error("test error message")

	if len(logger.buf) != 1 {
		t.Errorf("Expected 1 log entry, got %d", len(logger.buf))
	}
	if logger.buf[0].level != ErrorLevel {
		t.Errorf("Expected WarnLevel (due to bug), got %c", logger.buf[0].level)
	}
	if logger.buf[0].message != "test error message" {
		t.Errorf("Expected 'test error message', got '%s'", logger.buf[0].message)
	}
}

func TestRequestLogger_Errorf(t *testing.T) {
	var buf bytes.Buffer
	logger := &requestLogger{
		id:  "test-123",
		buf: make([]logEntry, 0),
		w:   &buf,
	}

	logger.Errorf("test error %v", errors.New("failed"))

	if len(logger.buf) != 1 {
		t.Errorf("Expected 1 log entry, got %d", len(logger.buf))
	}
	if logger.buf[0].level != ErrorLevel {
		t.Errorf("Expected ErrorLevel, got %c", logger.buf[0].level)
	}
	if logger.buf[0].message != "test error failed" {
		t.Errorf("Expected 'test error failed', got '%s'", logger.buf[0].message)
	}
}

func TestWithLogger(t *testing.T) {
	ctx := context.Background()
	newCtx := WithLogger(ctx)

	logger := FromContext(newCtx)

	if logger.id == "" {
		t.Error("Expected logger to have ID, got empty string")
	}
	if len(logger.buf) != 0 {
		t.Errorf("Expected empty buffer, got %d entries", len(logger.buf))
	}
}

func TestFromContext_WithLogger(t *testing.T) {
	ctx := context.Background()
	ctx = WithLogger(ctx)

	logger := FromContext(ctx)

	if logger.id == "" {
		t.Error("Expected logger to have ID, got empty string")
	}
}

func TestFromContext_NoLogger(t *testing.T) {
	ctx := context.Background()

	logger := FromContext(ctx)

	if logger.id != "noop" {
		t.Errorf("Expected 'noop' ID, got '%s'", logger.id)
	}

	if logger.w != io.Discard {
		t.Error("Expected io.Discard writer for noop logger")
	}
}

func TestRequestLogger_FlushIf_WithError(t *testing.T) {
	var buf bytes.Buffer
	logger := &requestLogger{
		id:  "test-123",
		buf: make([]logEntry, 0),
		w:   &buf,
	}

	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")

	testErr := errors.New("test error")
	logger.FlushIf(testErr)

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	if len(lines) != 4 {
		t.Errorf("Expected 4 lines of output, got %d", len(lines))
	}

	expectedLines := []string{
		"[test-123] D: debug message",
		"[test-123] I: info message",
		"[test-123] W: warn message",
		"[test-123] E: test error",
	}

	for i, expected := range expectedLines {
		if i < len(lines) && lines[i] != expected {
			t.Errorf("Line %d: expected '%s', got '%s'", i, expected, lines[i])
		}
	}
}

func TestRequestLogger_FlushIf_NoError(t *testing.T) {
	var buf bytes.Buffer
	logger := &requestLogger{
		id:  "test-123",
		buf: make([]logEntry, 0),
		w:   &buf,
	}

	logger.Debug("debug message")
	logger.Info("info message")

	logger.FlushIf(nil)

	output := buf.String()
	if output != "" {
		t.Errorf("Expected no output when flushing with nil error, got '%s'", output)
	}
}

func TestRequestLogger_Flush(t *testing.T) {
	var buf bytes.Buffer
	logger := &requestLogger{
		id:  "test-456",
		buf: make([]logEntry, 0),
		w:   &buf,
	}

	logger.Debug("debug message")
	logger.Info("info message")

	logger.Flush()

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	if len(lines) != 2 {
		t.Errorf("Expected 2 lines of output, got %d", len(lines))
	}

	expectedLines := []string{
		"[test-456] D: debug message",
		"[test-456] I: info message",
	}

	for i, expected := range expectedLines {
		if i < len(lines) && lines[i] != expected {
			t.Errorf("Line %d: expected '%s', got '%s'", i, expected, lines[i])
		}
	}
}

func TestRequestLogger_EmptyFlush(t *testing.T) {
	var buf bytes.Buffer
	logger := &requestLogger{
		id:  "test-789",
		buf: make([]logEntry, 0),
		w:   &buf,
	}

	logger.Flush()

	output := buf.String()
	if output != "" {
		t.Errorf("Expected no output when flushing empty buffer, got '%s'", output)
	}
}

func TestRequestLogger_MultipleLogLevels(t *testing.T) {
	var buf bytes.Buffer
	logger := &requestLogger{
		id:  "multi-test",
		buf: make([]logEntry, 0),
		w:   &buf,
	}

	logger.Debug("debug")
	logger.Debugf("debugf %d", 1)
	logger.Info("info")
	logger.Infof("infof %d", 2)
	logger.Warn("warn")
	logger.Warnf("warnf %d", 3)
	logger.Error("error")
	logger.Errorf("errorf %d", 4)

	if len(logger.buf) != 8 {
		t.Errorf("Expected 8 log entries, got %d", len(logger.buf))
	}

	expectedLevels := []Level{DebugLevel, DebugLevel, InfoLevel, InfoLevel, WarnLevel, WarnLevel, ErrorLevel, ErrorLevel}
	for i, expected := range expectedLevels {
		if i < len(logger.buf) && logger.buf[i].level != expected {
			t.Errorf("Entry %d: expected level %c, got %c", i, expected, logger.buf[i].level)
		}
	}
}

func TestPoolReuse(t *testing.T) {
	ctx1 := WithLogger(context.Background())
	logger1 := FromContext(ctx1)
	logger1.Debug("test message")

	id1 := logger1.id

	logger1.FlushIf(nil)

	ctx2 := WithLogger(context.Background())
	logger2 := FromContext(ctx2)

	if len(logger2.buf) != 0 {
		t.Errorf("Expected empty buffer from pool reuse, got %d entries", len(logger2.buf))
	}

	if id1 == logger2.id {
		t.Error("Expected different IDs for different logger instances")
	}
}

func TestConcurrentUsage(t *testing.T) {
	var wg sync.WaitGroup
	const numGoroutines = 100

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			ctx := WithLogger(context.Background())
			logger := FromContext(ctx)

			logger.Debug(fmt.Sprintf("debug message %d", id))
			logger.Info(fmt.Sprintf("info message %d", id))
			logger.Warn(fmt.Sprintf("warn message %d", id))

			if id%2 == 0 {
				logger.FlushIf(errors.New("test error"))
			} else {
				logger.FlushIf(nil)
			}
		}(i)
	}

	wg.Wait()
}

type failingWriter struct {
	failCount int
	callCount int
}

func (fw *failingWriter) Write(p []byte) (n int, err error) {
	fw.callCount++
	if fw.callCount <= fw.failCount {
		return 0, errors.New("write failed")
	}
	return len(p), nil
}

func TestRequestLogger_FlushWithWriteError(t *testing.T) {
	fw := &failingWriter{failCount: 10} // Fail all writes
	logger := &requestLogger{
		id:  "test-error",
		buf: make([]logEntry, 0),
		w:   fw,
	}

	logger.Debug("debug message")
	logger.Info("info message")

	logger.Flush()

	if fw.callCount != 2 {
		t.Errorf("Expected 2 write calls, got %d", fw.callCount)
	}
}

func TestRequestLogger_FlushIfWithWriteError(t *testing.T) {
	fw := &failingWriter{failCount: 10} // Fail all writes
	logger := &requestLogger{
		id:  "test-error",
		buf: make([]logEntry, 0),
		w:   fw,
	}

	logger.Debug("debug message")
	logger.Info("info message")

	// Should not panic even with write errors
	logger.FlushIf(errors.New("test error"))

	if fw.callCount != 3 { // 2 buffered entries + 1 error
		t.Errorf("Expected 3 write calls, got %d", fw.callCount)
	}
}

// BenchmarkRequestLogger_Debug benchmarks the Debug method
func BenchmarkRequestLogger_Debug(b *testing.B) {
	logger := &requestLogger{
		id:  "bench-test",
		buf: make([]logEntry, 0, 32),
		w:   io.Discard,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Debug("debug message")
	}
}

// BenchmarkRequestLogger_Debugf benchmarks the Debugf method
func BenchmarkRequestLogger_Debugf(b *testing.B) {
	logger := &requestLogger{
		id:  "bench-test",
		buf: make([]logEntry, 0, 32),
		w:   io.Discard,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Debugf("debug message %d", i)
	}
}

// BenchmarkRequestLogger_Info benchmarks the Info method
func BenchmarkRequestLogger_Info(b *testing.B) {
	logger := &requestLogger{
		id:  "bench-test",
		buf: make([]logEntry, 0, 32),
		w:   io.Discard,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("info message")
	}
}

// BenchmarkRequestLogger_Infof benchmarks the Infof method
func BenchmarkRequestLogger_Infof(b *testing.B) {
	logger := &requestLogger{
		id:  "bench-test",
		buf: make([]logEntry, 0, 32),
		w:   io.Discard,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Infof("info message %d", i)
	}
}

// BenchmarkRequestLogger_Warn benchmarks the Warn method
func BenchmarkRequestLogger_Warn(b *testing.B) {
	logger := &requestLogger{
		id:  "bench-test",
		buf: make([]logEntry, 0, 32),
		w:   io.Discard,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Warn("warn message")
	}
}

// BenchmarkRequestLogger_Warnf benchmarks the Warnf method
func BenchmarkRequestLogger_Warnf(b *testing.B) {
	logger := &requestLogger{
		id:  "bench-test",
		buf: make([]logEntry, 0, 32),
		w:   io.Discard,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Warnf("warn message %d", i)
	}
}

// BenchmarkRequestLogger_Error benchmarks the Error method
func BenchmarkRequestLogger_Error(b *testing.B) {
	logger := &requestLogger{
		id:  "bench-test",
		buf: make([]logEntry, 0, 32),
		w:   io.Discard,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Error("error message")
	}
}

// BenchmarkRequestLogger_Errorf benchmarks the Errorf method
func BenchmarkRequestLogger_Errorf(b *testing.B) {
	logger := &requestLogger{
		id:  "bench-test",
		buf: make([]logEntry, 0, 32),
		w:   io.Discard,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Errorf("error message %d", i)
	}
}

// BenchmarkRequestLogger_FlushIf_WithError benchmarks FlushIf with error
func BenchmarkRequestLogger_FlushIf_WithError(b *testing.B) {
	testErr := errors.New("test error")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger := &requestLogger{
			id:  "bench-test",
			buf: make([]logEntry, 0, 32),
			w:   io.Discard,
		}
		logger.Debug("debug message")
		logger.Info("info message")
		logger.Warn("warn message")
		logger.FlushIf(testErr)
	}
}

// BenchmarkRequestLogger_FlushIf_NoError benchmarks FlushIf without error
func BenchmarkRequestLogger_FlushIf_NoError(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger := &requestLogger{
			id:  "bench-test",
			buf: make([]logEntry, 0, 32),
			w:   io.Discard,
		}
		logger.Debug("debug message")
		logger.Info("info message")
		logger.Warn("warn message")
		logger.FlushIf(nil)
	}
}

// BenchmarkRequestLogger_Flush benchmarks Flush method
func BenchmarkRequestLogger_Flush(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger := &requestLogger{
			id:  "bench-test",
			buf: make([]logEntry, 0, 32),
			w:   io.Discard,
		}
		logger.Debug("debug message")
		logger.Info("info message")
		logger.Warn("warn message")
		logger.Flush()
	}
}

// BenchmarkWithLogger benchmarks context logger creation
func BenchmarkWithLogger(b *testing.B) {
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = WithLogger(ctx)
	}
}

// BenchmarkFromContext benchmarks logger retrieval from context
func BenchmarkFromContext(b *testing.B) {
	ctx := WithLogger(context.Background())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FromContext(ctx)
	}
}

// BenchmarkFromContext_NoLogger benchmarks fallback logger creation
func BenchmarkFromContext_NoLogger(b *testing.B) {
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FromContext(ctx)
	}
}

// BenchmarkFullWorkflow benchmarks a complete request lifecycle
func BenchmarkFullWorkflow(b *testing.B) {
	testErr := errors.New("test error")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := WithLogger(context.Background())
		logger := FromContext(ctx)

		logger.Debug("processing request")
		logger.Info("validating input")
		logger.Warn("potential issue")
		logger.Error("operation failed")

		if i%2 == 0 {
			logger.FlushIf(testErr)
		} else {
			logger.FlushIf(nil)
		}
	}
}

// BenchmarkFullWorkflow_WithoutError benchmarks successful request workflow
func BenchmarkFullWorkflow_WithoutError(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := WithLogger(context.Background())
		logger := FromContext(ctx)

		logger.Debug("processing request")
		logger.Info("validating input")
		logger.Warn("potential issue")

		logger.FlushIf(nil)
	}
}

// BenchmarkMemoryGrowth benchmarks buffer growth behavior
func BenchmarkMemoryGrowth(b *testing.B) {
	logger := &requestLogger{
		id:  "bench-test",
		buf: make([]logEntry, 0, 32),
		w:   io.Discard,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Debug("debug message")
		logger.Info("info message")
		logger.Warn("warn message")
		logger.Error("error message")

		// Simulate buffer growth beyond initial capacity
		if i%100 == 0 {
			logger.buf = logger.buf[:0] // Reset buffer
		}
	}
}

// BenchmarkConcurrentUsage benchmarks concurrent logger usage
func BenchmarkConcurrentUsage(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ctx := WithLogger(context.Background())
			logger := FromContext(ctx)

			logger.Debug("debug message")
			logger.Info("info message")
			logger.Warn("warn message")

			logger.FlushIf(nil)
		}
	})
}

// BenchmarkConcurrentUsage_WithError benchmarks concurrent logger usage with errors
func BenchmarkConcurrentUsage_WithError(b *testing.B) {
	testErr := errors.New("test error")

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ctx := WithLogger(context.Background())
			logger := FromContext(ctx)

			logger.Debug("debug message")
			logger.Info("info message")
			logger.Error("error message")

			logger.FlushIf(testErr)
		}
	})
}

// BenchmarkPoolReuse benchmarks pool reuse efficiency
func BenchmarkPoolReuse(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := WithLogger(context.Background())
		logger := FromContext(ctx)
		logger.Debug("test message")
		logger.FlushIf(nil) // Return to pool
	}
}

// BenchmarkStringFormatting benchmarks string formatting overhead
func BenchmarkStringFormatting(b *testing.B) {
	logger := &requestLogger{
		id:  "bench-test",
		buf: make([]logEntry, 0, 32),
		w:   io.Discard,
	}

	b.Run("Direct", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			logger.Debug("static message")
		}
	})

	b.Run("Formatted", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			logger.Debugf("message %d", i)
		}
	})

	b.Run("ComplexFormatted", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			logger.Debugf("complex message %d with %s and %v", i, "string", errors.New("error"))
		}
	})
}

// BenchmarkUUIDGeneration benchmarks UUID generation overhead
func BenchmarkUUIDGeneration(b *testing.B) {
	logger := &requestLogger{
		id:  "bench-test",
		buf: make([]logEntry, 0, 32),
		w:   io.Discard,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.reset()
	}
}

// BenchmarkCompareWithStandardLog benchmarks against standard log (if available)
func BenchmarkCompareWithStandardLog(b *testing.B) {
	b.Run("RequestLogger", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ctx := WithLogger(context.Background())
			logger := FromContext(ctx)
			logger.Info("test message")
			logger.FlushIf(nil)
		}
	})
}

// BenchmarkLargeBuffers benchmarks behavior with large log buffers
func BenchmarkLargeBuffers(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("Size_%d", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				logger := &requestLogger{
					id:  "bench-test",
					buf: make([]logEntry, 0, 32),
					w:   io.Discard,
				}

				// Fill buffer with entries
				for j := 0; j < size; j++ {
					logger.Debug("debug message")
				}

				logger.FlushIf(errors.New("test error"))
			}
		})
	}
}

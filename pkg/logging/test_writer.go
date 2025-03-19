package logging

import "strings"

// TestWriter is a simple io.Writer implementation that writes to a string builder.
// It is useful for testing purposes - to check what was written by the logger.
type TestWriter struct {
	builder strings.Builder
}

// Write satisfies the io.Writer interface.
func (w *TestWriter) Write(p []byte) (n int, err error) {
	return w.builder.Write(p) //nolint: wrapcheck // no need to wrap the error for testing
}

// String returns the content written to the writer.
func (w *TestWriter) String() string {
	return w.builder.String()
}

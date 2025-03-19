package logging

import "strings"

type TestWriter struct {
	builder strings.Builder
}

func (w *TestWriter) Write(p []byte) (n int, err error) {
	return w.builder.Write(p)
}

func (w *TestWriter) String() string {
	return w.builder.String()
}

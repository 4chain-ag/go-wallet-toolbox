package logging

import (
	"fmt"
	"log/slog"
	"os"
)

func Fatalf(logger *slog.Logger, err error, msg string, args ...any) {
	logger.Error("Fatal error: "+fmt.Sprintf(msg, args...), slog.String("error", err.Error()))
	os.Exit(1)
}

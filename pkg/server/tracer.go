package server

import (
	"context"
	"github.com/filecoin-project/go-jsonrpc"
	"log/slog"
	"reflect"
)

func tracer(logger *slog.Logger) jsonrpc.Tracer {
	return func(method string, _ []reflect.Value, _ []reflect.Value, err error) {
		level := slog.LevelInfo
		args := []slog.Attr{
			slog.String("method", method),
		}

		if err != nil {
			level = slog.LevelError
			args = append(args, slog.String("error", err.Error()))
		}

		logger.LogAttrs(context.Background(), level, "Handling RPC call", args...)
	}
}

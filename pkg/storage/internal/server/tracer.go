package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/logging"
	"github.com/filecoin-project/go-jsonrpc"
	"log/slog"
	"reflect"
)

func tracer(logger *slog.Logger) jsonrpc.Tracer {
	return func(method string, params []reflect.Value, results []reflect.Value, err error) {
		level := slog.LevelInfo
		args := []slog.Attr{
			slog.String("method", method),
		}

		if err != nil {
			level = slog.LevelError
			args = append(args, slog.String("error", err.Error()))
		}

		if logging.IsDebug(logger) {
			for i, param := range params {
				args = append(args, slog.Any(fmt.Sprintf("param_%d", i), reflectValueToLoggable(param)))
			}
			for i, result := range results {
				args = append(args, slog.Any(fmt.Sprintf("result_%d", i), reflectValueToLoggable(result)))
			}
		}

		logger.LogAttrs(context.Background(), level, "Handling RPC call", args...)
	}
}

func reflectValueToLoggable(v reflect.Value) string {
	if !v.IsValid() {
		return "<invalid>"
	}

	underlyingValue := v.Interface()

	if err, ok := underlyingValue.(error); ok {
		return fmt.Sprintf("<error: %v>", err)
	}

	if ctx, ok := underlyingValue.(context.Context); ok {
		return fmt.Sprintf("<context: %v>", ctx)
	}

	jsonBytes, err := json.Marshal(underlyingValue)
	if err != nil {
		return fmt.Sprintf("<error: %v>", err)
	}

	return string(jsonBytes)
}

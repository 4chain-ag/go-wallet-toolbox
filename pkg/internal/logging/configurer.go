package logging

import (
	"io"
	"log/slog"
)

type HandlerType string

const (
	JSONHandler HandlerType = "json"
	TextHandler HandlerType = "text"
)

type Configurer interface {
	WithLevel(level slog.Level) HandlerConfigurer
}

type HandlerConfigurer interface {
	WithCustomHandler(handler slog.Handler) LoggerMaker
	WithHandler(handlerType HandlerType, writer io.Writer) LoggerMaker
}

type LoggerMaker interface {
	Logger() *slog.Logger
}

type configurer struct {
	handler slog.Handler
	level   *slog.LevelVar
}

func New() Configurer {
	return &configurer{
		level: new(slog.LevelVar),
	}
}

func (c *configurer) WithLevel(level slog.Level) HandlerConfigurer {
	c.level.Set(level)
	return c
}

func (c *configurer) WithCustomHandler(handler slog.Handler) LoggerMaker {
	c.handler = handler
	return c
}

func (c *configurer) WithHandler(handlerType HandlerType, writer io.Writer) LoggerMaker {
	opts := &slog.HandlerOptions{Level: c.level}

	switch handlerType {
	case JSONHandler:
		c.handler = slog.NewJSONHandler(writer, opts)
	case TextHandler:
		c.handler = slog.NewTextHandler(writer, opts)
	default:
		panic("unsupported handler type")
	}
	return c
}

func (c *configurer) Logger() *slog.Logger {
	return slog.New(c.handler)
}

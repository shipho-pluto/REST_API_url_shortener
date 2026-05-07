package emptylog

import (
	"context"
	"log/slog"
)

func NewEmptyLogger() *slog.Logger {
	return slog.New(NewEmptyHandler())
}

type EmptyHandler struct{}

func NewEmptyHandler() *EmptyHandler {
	return &EmptyHandler{}
}

func (h *EmptyHandler) Handle(_ context.Context, _ slog.Record) error {
	return nil
}

func (h *EmptyHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

func (h *EmptyHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}

func (h *EmptyHandler) WithGroup(_ string) slog.Handler {
	return h
}

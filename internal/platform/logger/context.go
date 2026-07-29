package logger

import (
	"context"

	"go.uber.org/zap"
)

type ctxKey struct{}

func WithContext(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, l)
}
func FromContext(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
		return l
	}
	return global // fallback: si no hay logger en context, usa el global
}

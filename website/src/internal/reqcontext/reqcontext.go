package reqcontext

import (
	"context"

	"valette.software/internal/i18n"
)

const requestContextKey = iota

type ReqContext struct {
	Localizer   i18n.Localizer
	CurrentPath string
}

func SetValue(ctx context.Context, value ReqContext) context.Context {
	return context.WithValue(ctx, requestContextKey, value)
}

func GetValue(ctx context.Context) ReqContext {
	value, ok := ctx.Value(requestContextKey).(ReqContext)

	if !ok {
		return ReqContext{}
	}

	return value
}

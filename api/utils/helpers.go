package utils

import (
	"context"
)

type tokenKey struct{}

func CreateContextWithToken(ctx context.Context, token interface{}) context.Context {
	return context.WithValue(ctx, tokenKey{}, token)
}

func GetTokenFromContext(ctx context.Context) interface{} {
	return ctx.Value(tokenKey{})
}

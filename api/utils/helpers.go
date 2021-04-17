package utils

import (
	"context"
)

type tokenKey struct{}
type sessionKey struct{}

type SessionData struct {
	UserID string
	Origin string
}

func CreateContextWithToken(ctx context.Context, token interface{}) context.Context {
	return context.WithValue(ctx, tokenKey{}, token)
}

func GetTokenFromContext(ctx context.Context) interface{} {
	return ctx.Value(tokenKey{})
}

func CreateContextWithSessionData(ctx context.Context, sessionID, origin string) context.Context {
	return context.WithValue(ctx, sessionKey{}, &SessionData{UserID: sessionID, Origin: origin})
}

func GetSessionDataFromCtx(ctx context.Context) *SessionData {
	session, ok := ctx.Value(sessionKey{}).(*SessionData)
	if !ok {
		return &SessionData{}
	}

	return session
}

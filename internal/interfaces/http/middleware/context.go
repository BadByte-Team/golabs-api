package middleware

import "context"

type UserContext struct {
	UserID string
	Role   string
	Banned bool
}

type ctxKey string

const userKey ctxKey = "user"

func WithUser(ctx context.Context, user UserContext) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func GetUser(ctx context.Context) (UserContext, bool) {
	user, ok := ctx.Value(userKey).(UserContext)
	return user, ok
}

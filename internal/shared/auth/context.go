package auth

import "context"

type Context struct {
	UserID int
	Role   string
}

type contextKey struct{}

func WithContext(ctx context.Context, authCtx *Context) context.Context {
	return context.WithValue(ctx, contextKey{}, authCtx)
}

func FromContext(ctx context.Context) (*Context, bool) {
	authCtx, ok := ctx.Value(contextKey{}).(*Context)
	return authCtx, ok
}

func MustFromContext(ctx context.Context) *Context {
	authCtx, ok := FromContext(ctx)
	if !ok {
		panic("missing auth context")
	}
	return authCtx
}

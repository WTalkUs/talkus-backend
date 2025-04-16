package middleware

type ContextKey string

const (
	AuthUserKey ContextKey = "authenticated_user"
)
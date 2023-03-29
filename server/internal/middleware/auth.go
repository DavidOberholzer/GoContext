package middleware

import (
	"context"
	"net/http"
	"strings"
)

type ctxKey int

const userCtxKey ctxKey = iota

func JSONResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// TokenAuth checks access token from super secret access thingy, and then puts the found userID on context.
func TokenAuth(next http.Handler) http.Handler {
	hypotheticalAccessThingy := map[string]string{
		"1234": "dave",
		"7890": "jen",
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authHeader := r.Header.Get("Authorization")
		parts := strings.Split(authHeader, "Token ")
		token := ""
		if len(parts) == 2 {
			token = parts[1]
		}

		userID, _ := hypotheticalAccessThingy[token]

		next.ServeHTTP(w, r.WithContext(context.WithValue(ctx, userCtxKey, userID)))
	})
}

func GetUserIDFromCtx(ctx context.Context) string {
	value, _ := ctx.Value(userCtxKey).(string)
	return value
}

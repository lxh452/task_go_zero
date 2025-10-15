package middleware

import (
    "net/http"

    "github.com/zeromicro/go-zero/rest/httpx"
    "github.com/zeromicro/go-zero/rest/token"
)

// JwtMiddleware validates JWT on protected routes.
type JwtMiddleware struct {
    verifier *token.TokenVerifier
}

func NewJwtMiddleware(secret string) *JwtMiddleware {
    return &JwtMiddleware{verifier: token.NewTokenVerifier(secret)}
}

func (m *JwtMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if err := m.verifier.VerifyToken(r.Header.Get("Authorization")); err != nil {
            httpx.Error(w, err)
            return
        }
        next(w, r)
    }
}



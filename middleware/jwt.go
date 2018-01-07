package middleware

import (
	"rady"
	"github.com/labstack/echo/middleware"
)

const (
	HeaderAuthorization = "Authorization"
	AlgorithmHS256      = "HS256"
	AuthScheme          = "Bearer"
	ContextKey          = "user"
)

type JWTMiddleware struct {
	rady.Middleware
	SigningKey    *string `value:"rady.jwt.signing_key"`
	TokenLookup   *string `value:"rady.jwt.token_lookup" default:"header:Authorization"`
	ContextKey    *string `value:"rady.jwt.context_key" default:"user"`
	SigningMethod *string `value:"rady.jwt.signing_method" default:"HS256"`
	AuthScheme    *string `value:"rady.jwt.auth_scheme" default:"Bearer"`
}

func (j *JWTMiddleware) GetWithConfig(next rady.HandlerFunc) rady.HandlerFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    []byte(*j.SigningKey),
		TokenLookup:   *j.TokenLookup,
		ContextKey:    *j.ContextKey,
		SigningMethod: *j.SigningMethod,
		AuthScheme:    *j.AuthScheme,
	})(next)
}

package middleware

import (
	"github.com/labstack/echo/middleware"
	"rady"
)

type Logger struct {
	rady.Middleware
}

func Log(next rady.HandlerFunc) rady.HandlerFunc {
	return middleware.Logger()(next)
}

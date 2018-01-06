package middleware

import (
	"rady"
	"github.com/labstack/echo/middleware"
)

type Logger struct {
	rady.Middleware
}

func Log(next rady.HandlerFunc) rady.HandlerFunc {
	return middleware.Logger()(next)
}
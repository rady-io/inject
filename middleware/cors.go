package middleware

import (
	"rady"
	"github.com/labstack/echo/middleware"
)

const (
	CONNECT = "CONNECT"
	DELETE  = "DELETE"
	GET     = "GET"
	HEAD    = "HEAD"
	OPTIONS = "OPTIONS"
	PATCH   = "PATCH"
	POST    = "POST"
	PUT     = "PUT"
	TRACE   = "TRACE"
)

var (
	AllowOrigins = []interface{}{"*"}
	AllowMethods = []interface{}{GET, HEAD, PUT, PATCH, POST, DELETE}
)

type CORSMiddleware struct {
	rady.Middleware
	// AllowOrigin defines a list of origins that may access the resource.
	// Optional. Default value *[]interface{}{"*"}.
	AllowOrigins *[]interface{} `value:"allow_origins" default:"[\"*\"]"`

	// AllowMethods defines a list methods allowed when accessing the resource.
	// This is used in response to a preflight request.
	// Optional. Default value DefaultCORSConfig.AllowMethods.

	AllowMethods *[]interface{} `value:"allow_methods" default:"[\"GET\", \"HEAD\", \"PUT\", \"POST\", \"PATCH\", \"DELETE\"]"`

	// AllowHeaders defines a list of request headers that can be used when
	// making the actual request. This in response to a preflight request.
	// Optional. Default value *[]interface{}{}.
	AllowHeaders *[]interface{} `value:"allow_headers" default:"[]"`

	// AllowCredentials indicates whether or not the response to the request
	// can be exposed when the credentials flag is true. When used as part of
	// a response to a preflight request, this indicates whether or not the
	// actual request can be made using credentials.
	// Optional. Default value false.
	AllowCredentials *bool `value:"allow_credentials" default:"false"`

	// ExposeHeaders defines a whitelist headers that clients are allowed to
	// access.
	// Optional. Default value *[]interface{}{}.
	ExposeHeaders *[]interface{} `value:"expose_headers" default:"[]"`

	// MaxAge indicates how long (in seconds) the results of a preflight request
	// can be cached.
	// Optional. Default value 0.
	MaxAge *int `value:"max_age" default:"0"`

	App *rady.Application
}

func InterfaceToString(interfaces *[]interface{}) []string {
	stringSlice := make([]string, 0)
	for _, value := range *interfaces {
		str, ok := value.(string)
		if !ok {
			continue
		}
		stringSlice = append(stringSlice, str)
	}
	return stringSlice
}

func (cors *CORSMiddleware) MiddlewareWithConfig(next rady.HandlerFunc) rady.HandlerFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     InterfaceToString(cors.AllowOrigins),
		AllowMethods:     InterfaceToString(cors.AllowMethods),
		AllowHeaders:     InterfaceToString(cors.AllowHeaders),
		AllowCredentials: *cors.AllowCredentials,
		ExposeHeaders:    InterfaceToString(cors.ExposeHeaders),
		MaxAge:           *cors.MaxAge,
	})(next)
}

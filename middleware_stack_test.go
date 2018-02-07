package rady

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

type (
	LoggerMiddleware struct {
		Middleware
	}

	CORS struct {
		Middleware
	}

	JWT struct {
		Middleware
	}

	RootRouter struct {
		Router `prefix:"/api"`
		*LoggerMiddleware
		*CORS
		*AuthRouter
	}

	AuthRouter struct {
		Router `prefix:"/auth"`
		*JWT
	}

	MidStackRoot struct {
		*RootRouter
	}

	MidStackTest struct {
		Testing
		App *Application
	}
)

func (logger *LoggerMiddleware) DefaultLogger(next HandlerFunc) HandlerFunc {
	return func(context Context) error {
		return nil
	}
}

func (cors *CORS) DefaultCORS(next HandlerFunc) HandlerFunc {
	return func(context Context) error {
		return nil
	}
}

func (jwt *JWT) DefaultJWT(next HandlerFunc) HandlerFunc {
	return func(context Context) error {
		return nil
	}
}

func (mid *MidStackTest) TestMidStacks(t *testing.T) {
	APIStack, ok := mid.App.MiddlewareStackMap["/api"]
	assert.True(t, ok)
	assert.Equal(t, "DefaultCORS", APIStack.Stack[0].Name)
	assert.Equal(t, "DefaultLogger", APIStack.Stack[1].Name)

	AuthStack, ok := mid.App.MiddlewareStackMap["/api/auth"]
	assert.True(t, ok)
	assert.Equal(t, "DefaultJWT", AuthStack.Stack[0].Name)
}

func TestMidStack(t *testing.T) {
	CreateApplication(new(MidStackRoot)).PrepareTest().AddTest(new(MidStackTest)).Test(t)
}

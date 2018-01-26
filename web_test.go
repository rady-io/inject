package rady

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/labstack/echo/middleware"
)

func TestMiddlewareStack(t *testing.T) {
	cors := middleware.CORS()
	logger := middleware.Logger()
	corsName := "middleware.CORS"
	loggerName := "middleware.Logger"
	stack := NewMiddlewareStack()
	stack = stack.Push(NewMiddlewareContainer(corsName, cors))
	newStack := stack.Push(NewMiddlewareContainer(loggerName, logger))

	assert.Equal(t, len(stack.Stack), 1)
	assert.Equal(t, len(newStack.Stack), 2)
	assert.Equal(t, stack.Stack[0], newStack.Stack[1])
	assert.Equal(t, newStack.Stack[0].Name, loggerName)
	assert.Equal(t, newStack.Stack[1].Name, corsName)

	newStack = stack.PushStack(newStack)
	assert.Equal(t, len(newStack.Stack), 3)
	assert.Equal(t, newStack.Stack[1], newStack.Stack[2])
	assert.Equal(t, newStack.Stack[0].Name, loggerName)
	assert.Equal(t, newStack.Stack[1].Name, corsName)
	assert.Equal(t, newStack.Stack[2].Name, corsName)
}

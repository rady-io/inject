package rady

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/labstack/echo/middleware"
	"reflect"
	"fmt"
)

func TestMiddlewareStack(t *testing.T) {
	cors := middleware.CORS()
	logger := middleware.Logger()
	stack := NewMiddlewareStack()
	stack = stack.Push(cors)
	newStack := stack.Push(logger)

	assert.Equal(t, len(stack.Stack), 1)
	assert.Equal(t, len(newStack.Stack), 2)
	assert.Equal(t, reflect.ValueOf(stack.Stack[0]), reflect.ValueOf(newStack.Stack[1]))

	fmt.Printf("%#v\n", stack)
	fmt.Printf("%#v\n", newStack)

	newStack = stack.PushStack(newStack)
	assert.Equal(t, len(newStack.Stack), 3)
	assert.Equal(t, reflect.ValueOf(newStack.Stack[1]), reflect.ValueOf(newStack.Stack[2]))
	fmt.Printf("%#v\n", newStack)
}

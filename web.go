package rady

import "github.com/labstack/echo"

type (
	Context = echo.Context

	HandlerFunc = echo.HandlerFunc

	MiddlewareFunc = echo.MiddlewareFunc

	Group = echo.Group

	MiddlewareStack struct {
		Stack []MiddlewareFunc
	}
)

func NewMiddlewareStack() *MiddlewareStack {
	return &MiddlewareStack{
		Stack: make([]MiddlewareFunc, 0),
	}
}

func (stack *MiddlewareStack) Push(middleware MiddlewareFunc) *MiddlewareStack {
	newStack := make([]MiddlewareFunc, len(stack.Stack)+1)
	newStack[0] = middleware
	copy(newStack[1:], stack.Stack)
	return &MiddlewareStack{
		Stack: newStack,
	}
}

func (stack *MiddlewareStack) PushStack(frontStack *MiddlewareStack) *MiddlewareStack {
	newStack := make([]MiddlewareFunc, len(stack.Stack)+len(frontStack.Stack))
	copy(newStack, frontStack.Stack)
	copy(newStack[len(frontStack.Stack):], stack.Stack)
	return &MiddlewareStack{
		Stack: newStack,
	}
}

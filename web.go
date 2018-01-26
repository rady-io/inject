package rady

import "github.com/labstack/echo"

type (
	Context = echo.Context

	HandlerFunc = echo.HandlerFunc

	MiddlewareFunc = echo.MiddlewareFunc

	Group = echo.Group

	MiddlewareContainer struct {
		Name string
		Func MiddlewareFunc
	}

	MiddlewareStack struct {
		Stack []*MiddlewareContainer
	}
)

func NewMiddlewareContainer(Name string, Func MiddlewareFunc) *MiddlewareContainer {
	return &MiddlewareContainer{Name, Func}
}

func NewMiddlewareStack() *MiddlewareStack {
	return &MiddlewareStack{
		Stack: make([]*MiddlewareContainer, 0),
	}
}

func (stack *MiddlewareStack) Push(middleware *MiddlewareContainer) *MiddlewareStack {
	newStack := make([]*MiddlewareContainer, len(stack.Stack)+1)
	newStack[0] = middleware
	copy(newStack[1:], stack.Stack)
	return &MiddlewareStack{
		Stack: newStack,
	}
}

func (stack *MiddlewareStack) PushStack(frontStack *MiddlewareStack) *MiddlewareStack {
	newStack := make([]*MiddlewareContainer, len(stack.Stack)+len(frontStack.Stack))
	copy(newStack, frontStack.Stack)
	copy(newStack[len(frontStack.Stack):], stack.Stack)
	return &MiddlewareStack{
		Stack: newStack,
	}
}

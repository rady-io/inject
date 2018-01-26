package rady

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"reflect"
	"strings"
	"github.com/labstack/echo/middleware"
	"fmt"
)

type ComponentHandler struct {
	Component
}

func TestContainsField(t *testing.T) {
	handler := ComponentHandler{}
	assert.True(t, ContainsField(reflect.TypeOf(handler), Component{}), "Handler should contain Component")

	assert.True(t, ContainsField(reflect.TypeOf(&handler).Elem(), Component{}), "Handler should contain Component")

	assert.False(t, ContainsField(reflect.TypeOf(handler), Router{}), "Handler should not contain Router")
}

func TestContainsFields(t *testing.T) {
	handler := ComponentHandler{}
	assert.True(t, ContainsFields(reflect.TypeOf(handler), ComponentTypes), "Handler should contain some field in COMPONENT_TYPES")

	typesSet := make(map[reflect.Type]bool)
	typesSet[reflect.TypeOf(Service{})] = true
	assert.False(t, ContainsFields(reflect.TypeOf(handler), typesSet), "Handler should not contain some field in typesSet")
}

func TestGetNewPrefix(t *testing.T) {
	prefix := "/api/v1/"
	path := "/vote/:id/"
	assert.Equal(t, GetNewPrefix(prefix, path), "/api/v1/vote/:id")
}

func TestParseHandlerName(t *testing.T) {
	ok, _, _ := ParseHandlerName("Gett_my_name")
	assert.False(t, ok)

	ok, method, path := ParseHandlerName("Post")
	assert.True(t, ok)
	assert.Equal(t, POST{}, method)
	assert.Equal(t, "", path)

	ok, method, path = ParseHandlerName("GetMyName")
	assert.True(t, ok)
	assert.Equal(t, GET{}, method)
	assert.Equal(t, "my/name", path)

	ok, method, path = ParseHandlerName("GetUserUUID")
	assert.True(t, ok)
	assert.Equal(t, GET{}, method)
	assert.Equal(t, "user/:uuid", path)

	ok, method, path = ParseHandlerName("GetUserUUIDGo")
	assert.True(t, ok)
	assert.Equal(t, GET{}, method)
	assert.Equal(t, "user/:uuid/go", path)
}

func TestSplitByUpper(t *testing.T) {
	result := SplitByUpper("GetMyName")
	assert.Equal(t, strings.Join(result, "/"), "Get/My/Name")

	result = SplitByUpper("GetMyNameID")
	assert.Equal(t, strings.Join(result, "/"), "Get/My/Name/ID")

	result = SplitByUpper("GetMyNameIDSid")
	assert.Equal(t, strings.Join(result, "/"), "Get/My/Name/ID/Sid")

	result = SplitByUpper("Get")
	assert.Equal(t, strings.Join(result, "/"), "Get")
}

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
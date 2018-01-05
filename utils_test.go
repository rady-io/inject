package rady

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"reflect"
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

	ok, method, path := ParseHandlerName("GetMyName")
	assert.True(t, ok)
	assert.Equal(t, GET{}, method)
	assert.Equal(t, "", path)

	ok, method, path = ParseHandlerName("Put_name_ID")
	assert.True(t, ok)
	assert.Equal(t, PUT{}, method)
	assert.Equal(t, "name/:id", path)
}
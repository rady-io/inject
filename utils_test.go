package rhapsody

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
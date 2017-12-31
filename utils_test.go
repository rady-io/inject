package summer

import (
	"summer/types"
	"testing"
	"github.com/stretchr/testify/assert"
	"reflect"
)

type Handler struct {
	types.Component
}

func TestContainsField(t *testing.T) {
	handler := Handler{}
	assert.True(t, ContainsField(reflect.TypeOf(handler), types.Component{}), "Handler should contain Component")

	assert.False(t, ContainsField(reflect.TypeOf(handler), types.Router{}), "Handler should not contain Router")
}

func TestContainsFields(t *testing.T) {
	handler := Handler{}
	assert.True(t, ContainsFields(reflect.TypeOf(handler), types.COMPONENT_TYPES), "Handler should contain some field in COMPONENT_TYPES")

	typesSet := make(map[reflect.Type]bool)
	typesSet[reflect.TypeOf(types.Service{})] = true
	assert.False(t, ContainsFields(reflect.TypeOf(handler), typesSet), "Handler should not contain some field in typesSet")
}
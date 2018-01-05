package rady

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"reflect"
)

func TestComponents(t *testing.T)  {
	_, ok := COMPONENTS[COMPONENT]
	assert.True(t, ok, "COMPONENTS[COMPONENT] should be true")

	_, ok = ComponentTypes[reflect.TypeOf(Component{})]
	assert.True(t, ok, "COMPONENT_TYPES[reflect.TypeOf(Component{})] should be true")
}

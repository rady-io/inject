package rady

import "reflect"

// Testing is a tag to mark a struct as a Testing
type (
	Testing struct {
	}

	TestingBean struct {
		Type  reflect.Type
		Value reflect.Value
	}
)

func NewTestingBean(Type reflect.Type, Value reflect.Value) *TestingBean {
	return &TestingBean{Type, Value}
}

const TESTING = "testing"

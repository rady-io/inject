package bean

import "reflect"

type Bean struct {
	Type  reflect.Type
	Tag   reflect.StructTag
	Value reflect.Value
}

type Method struct {
	Value reflect.Value
	In   reflect.Type
	Name string
}

func NewBean(Value reflect.Value, Tag reflect.StructTag) *Bean {
	return &Bean{
		Type:  Value.Type(),
		Tag:   Tag,
		Value: Value,
	}
}

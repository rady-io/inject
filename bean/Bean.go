package bean

import "reflect"

type Bean struct {
	//Type  reflect.Type
	Tag   reflect.StructTag
	Value reflect.Value
}

type Method struct {
	Value reflect.Value
	In   reflect.Type
	Name string
}

type Controller struct {
	Value reflect.Value
	Tag reflect.StructTag
}

type Middleware struct {
	Value reflect.Value
	Tag reflect.StructTag
}

func NewBean(Value reflect.Value, Tag reflect.StructTag) *Bean {
	return &Bean{
		//Type:  Value.Type(),
		Tag:   Tag,
		Value: Value,
	}
}

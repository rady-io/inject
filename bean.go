package rhapsody

import (
	"reflect"
	"github.com/tidwall/gjson"
)

type Bean struct {
	//Type  reflect.Type
	Tag   reflect.StructTag
	Value reflect.Value
}

type Method struct {
	Value reflect.Value
	Ins   []reflect.Type
	Name  string
}

type ParamBean struct {
	Value      reflect.Value
	MethodBean *Method
}

type ValueBean struct {
	Value      gjson.Result
	ValueMap   map[reflect.Type]reflect.Value
	ParamSlice []*ParamBean
}

type CtrlBean struct {
	Value reflect.Value
	Tag   reflect.StructTag
}

type MdWareBean struct {
	Value reflect.Value
	Tag   reflect.StructTag
}

func NewBean(Value reflect.Value, Tag reflect.StructTag) *Bean {
	return &Bean{
		//Type:  Value.Type(),
		Tag:   Tag,
		Value: Value,
	}
}

func NewBeanMethod(Value reflect.Value, Name string) *Method {
	return &Method{
		Value: Value,
		Ins:   make([]reflect.Type, 0),
		Name:  Name,
	}
}

func NewParamBean(Value reflect.Value, MethodBean *Method) *ParamBean {
	return &ParamBean{
		Value:      Value,
		MethodBean: MethodBean,
	}
}

func NewValueBean(Value gjson.Result) *ValueBean  {
	return &ValueBean{
		Value:Value,
		ValueMap: make(map[reflect.Type]reflect.Value),
		ParamSlice: make([]*ParamBean, 0),
	}
}

func NewCtrlBean(Value reflect.Value, Tag reflect.StructTag) *CtrlBean {
	return &CtrlBean{
		//Type:  Value.Type(),
		Tag:   Tag,
		Value: Value,
	}
}

func NewMdWareBean(Value reflect.Value, Tag reflect.StructTag) *MdWareBean {
	return &MdWareBean{
		//Type:  Value.Type(),
		Tag:   Tag,
		Value: Value,
	}
}

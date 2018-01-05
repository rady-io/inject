package rady

import (
	"reflect"
	"github.com/tidwall/gjson"
)

// Bean contains the value and tag of a type
type Bean struct {
	Tag   reflect.StructTag
	Value reflect.Value
}

// Method contains value, param list and name of a 'BeanMethod'
type Method struct {
	Value    reflect.Value
	Ins      []reflect.Type
	Name     string
	OutValue reflect.Value
}

///*
//ParamBean contains value of a bean which is parameter of a bean method
//
//contains also the bean of method it belongs to
// */
//type ParamBean struct {
//	Value      reflect.Value
//	MethodBean *Method
//}

/*
ValueBean contains value from config file parsed by 'gjson'

ValueMap is different types the value converted to

ParamSlice is the param list contain this value
 */
type ValueBean struct {
	Value       gjson.Result
	ValueMap    map[reflect.Type]reflect.Value
	MethodSlice []*Method
}

/*
CtrlBean contains value and tag of a controller
 */
type CtrlBean struct {
	Value reflect.Value
	Tag   reflect.StructTag
}

/*
MdWareBean contains value and tag of a middleware
 */
type MdWareBean struct {
	Value reflect.Value
	Tag   reflect.StructTag
}

/*
NewBean is factory function of Bean
*/
func NewBean(Value reflect.Value, Tag reflect.StructTag) *Bean {
	return &Bean{
		Tag:   Tag,
		Value: Value,
	}
}

/*
NewBeanMethod is factory function of Method
*/
func NewBeanMethod(Value reflect.Value, Name string) *Method {
	return &Method{
		Value: Value,
		Ins:   make([]reflect.Type, 0),
		Name:  Name,
	}
}

///*
//NewParamBean is factory function of ParamBean
// */
//func NewParamBean(Value reflect.Value, MethodBean *Method) *ParamBean {
//	return &ParamBean{
//		Value:      Value,
//		MethodBean: MethodBean,
//	}
//}

/*
NewValueBean is factory function of ValueBean
 */
func NewValueBean(Value gjson.Result) *ValueBean {
	return &ValueBean{
		Value:       Value,
		ValueMap:    make(map[reflect.Type]reflect.Value),
		MethodSlice: make([]*Method, 0),
	}
}

/*
NewCtrlBean is factory function of CtrlBean
 */
func NewCtrlBean(Value reflect.Value, Tag reflect.StructTag) *CtrlBean {
	return &CtrlBean{
		//Type:  Value.Type(),
		Tag:   Tag,
		Value: Value,
	}
}

/*
NewMdWareBean is factory function of MdwareBean
 */
func NewMdWareBean(Value reflect.Value, Tag reflect.StructTag) *MdWareBean {
	return &MdWareBean{
		//Type:  Value.Type(),
		Tag:   Tag,
		Value: Value,
	}
}

package rady

import (
	"reflect"
	"github.com/tidwall/gjson"
	"os"
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
	InValues []reflect.Value
}

func (m *Method) LoadIns(app *Application) {
	for _, inType := range m.Ins {
		if ConfirmSameTypeInMap(app.BeanMap, inType) {
			if len(app.BeanMap[inType]) > 1 {
				app.Logger.Critical("There are more than one %s, please named it.", inType)
				os.Exit(1)
			}
			for _, bean := range app.BeanMap[inType] {
				m.InValues = append(m.InValues, bean.Value)
			}
		} else {
			newValue := reflect.New(inType.Elem()).Elem()
			app.load(inType, newValue, GetTagFromName(""))
			m.InValues = append(m.InValues, newValue)
		}
	}
}

func (m *Method) Call(app *Application) {
	params := make([]reflect.Value, 0)
	for _, value := range m.InValues {
		params = append(params, value.Addr())
	}
	result := m.Value.Call(params)
	if len(result) != 1 {
		app.Logger.Error("Result of %s is not a Component!!!", m.Name)
		os.Exit(1)
	}
	app.Logger.Debug("Result of %s set %s", m.Name, result[0].Elem())
	m.OutValue.Set(result[0].Elem())
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
	Value     gjson.Result
	ValueMap  map[reflect.Type]reflect.Value
	MethodSet map[*Method]bool
	Key       string
	Default   gjson.Result
}

func (v *ValueBean) Reload(a *Application) {
	newResult := gjson.Get(a.ConfigFile, v.Key)
	if newResult != v.Value {
		if !newResult.Exists() {
			a.Logger.Info("Key %s doesn't exist, use default value %s", v.Key, v.Default.String())
			newResult = v.Default
		}
		v.Value = newResult
		a.Logger.Debug("Reset Value '%s' to %s", v.Key, v.Default.String())
		v.resetValue()
		v.recallFactory(a)
	}
}

func (v *ValueBean) resetValue() {
	for Type, Value := range v.ValueMap {
		switch Type {
		case IntPtrType:
			Value.SetInt(v.Value.Int())
		case UintPtrType:
			Value.SetUint(v.Value.Uint())
		case FloatPtrType:
			Value.SetFloat(v.Value.Float())
		case StringPtrType:
			Value.SetString(v.Value.String())
		case BoolPtrType:
			Value.SetBool(v.Value.Bool())
		case TimePtrType:
			Value.Set(reflect.ValueOf(v.Value.Time()))
		case ArrayPtrType:
			Value.Set(reflect.ValueOf(v.Value.Array()))
		case MapPtrType:
			Value.Set(reflect.ValueOf(v.Value.Map()))
		}
	}
}

func (v *ValueBean) recallFactory(a *Application) {
	for Method := range v.MethodSet {
		if _, ok := a.FactoryToRecall[Method]; !ok {
			a.FactoryToRecall[Method] = true
		}
	}
}

func (v *ValueBean) SetValue(value reflect.Value, Type reflect.Type) bool {
	confValue, ok := v.ValueMap[Type]
	if ok {
		value.Set(confValue.Addr())
		return true
	}
	switch Type {
	case IntPtrType:
		result := v.Value.Int()
		v.ValueMap[Type] = reflect.ValueOf(&result).Elem()
	case UintPtrType:
		result := v.Value.Uint()
		v.ValueMap[Type] = reflect.ValueOf(&result).Elem()
	case FloatPtrType:
		result := v.Value.Float()
		v.ValueMap[Type] = reflect.ValueOf(&result).Elem()
	case StringPtrType:
		result := v.Value.String()
		v.ValueMap[Type] = reflect.ValueOf(&result).Elem()
	case BoolPtrType:
		result := v.Value.Bool()
		v.ValueMap[Type] = reflect.ValueOf(&result).Elem()
	case TimePtrType:
		result := v.Value.Time()
		v.ValueMap[Type] = reflect.ValueOf(&result).Elem()
	case ArrayPtrType:
		result := v.Value.Array()
		v.ValueMap[Type] = reflect.ValueOf(&result).Elem()
	case MapPtrType:
		result := v.Value.Map()
		v.ValueMap[Type] = reflect.ValueOf(&result).Elem()
	}

	confValue, ok = v.ValueMap[Type]
	if ok {
		value.Set(confValue.Addr())
		return true
	}

	return false
}

/*
CtrlBean contains value and tag of a controller
 */
type CtrlBean struct {
	Name  string
	Value reflect.Value
	Tag   reflect.StructTag
}

/*
MdWareBean contains value and tag of a middleware
 */
type MdWareBean struct {
	Name  string
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
		Value:    Value,
		Ins:      make([]reflect.Type, 0),
		InValues: make([]reflect.Value, 0),
		Name:     Name,
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
func NewValueBean(Value gjson.Result, key string, defaultValue gjson.Result) *ValueBean {
	return &ValueBean{
		Value:     Value,
		ValueMap:  make(map[reflect.Type]reflect.Value),
		MethodSet: make(map[*Method]bool),
		Key:       key,
		Default:   defaultValue,
	}
}

/*
NewCtrlBean is factory function of CtrlBean
 */
func NewCtrlBean(Value reflect.Value, Tag reflect.StructTag, Name string) *CtrlBean {
	return &CtrlBean{
		Name:  Name,
		Tag:   Tag,
		Value: Value,
	}
}

/*
NewMdWareBean is factory function of MdwareBean
 */
func NewMdWareBean(Value reflect.Value, Tag reflect.StructTag, Name string) *MdWareBean {
	return &MdWareBean{
		Name:  Name,
		Tag:   Tag,
		Value: Value,
	}
}

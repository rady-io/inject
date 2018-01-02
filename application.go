package rhapsody

import (
	"reflect"
	"github.com/labstack/echo"
	"fmt"
	"strings"
)

type Application struct {
	app             interface{}
	BeanMap         map[reflect.Type]map[string]*Bean
	BeanMethodMap   map[reflect.Type]map[string]*Method
	ValueBeanMap    map[string]*ValueBean
	CtrlBeanSlice   []*CtrlBean
	MdWareBeanSlice []*MdWareBean
	Server          *echo.Echo
	Logger          *Logger
	ConfigFile      string
}

func CreateApplication(app interface{}) *Application {
	if CheckFieldPtr(reflect.TypeOf(app)) {
		return (&Application{
			app:             app,
			BeanMap:         make(map[reflect.Type]map[string]*Bean),
			BeanMethodMap:   make(map[reflect.Type]map[string]*Method),
			ValueBeanMap:    make(map[string]*ValueBean),
			CtrlBeanSlice:   make([]*CtrlBean, 0),
			MdWareBeanSlice: make([]*MdWareBean, 0),
			Server:          echo.New(),
			Logger:          NewLogger(),
		}).init()
	}
	NewLogger().Errorf("%s is not kind of Ptr!!!\n", reflect.TypeOf(app).Name)
	return new(Application)
}

func (a *Application) init() *Application {
	return a.loadElem(a.Logger, *new(reflect.StructTag)).loadElem(a, *new(reflect.StructTag)).loadConfigFile()
}

func (a *Application) loadElem(elem interface{}, tag reflect.StructTag) *Application {
	Value := reflect.ValueOf(elem)
	Type := reflect.TypeOf(elem)
	return a.load(Type, Value.Elem(), tag)
}

func (a *Application) load(fieldType reflect.Type, Value reflect.Value, tag reflect.StructTag) *Application {
	name := GetBeanName(fieldType, tag)
	a.Logger.Debug("%s -> %s", name, Value.Type())
	if ConfirmAddBeanMap(a.BeanMap, fieldType, name) {
		newBean := NewBean(Value, tag)
		a.BeanMap[fieldType][name] = newBean
	}
	return a
}

func (a *Application) Run() {
	app := a.app
	appType := reflect.TypeOf(app).Elem()
	for i := 0; i < appType.NumField(); i++ {
		field := appType.Field(i)
		if CheckConfiguration(field) {
			fieldValue := reflect.New(field.Type.Elem()).Elem() // save Elem in Bean
			a.loadField(field, fieldValue)
			for i := 0; i < fieldValue.Addr().NumMethod(); i++ {
				name := fieldValue.Addr().Type().Method(i).Name
				a.loadBeanMethod(fieldValue.Addr().MethodByName(name), name)
			}
		}
	}
}

func (a *Application) loadField(Field reflect.StructField, FieldValue reflect.Value) {
	fieldType := Field.Type
	a.load(fieldType, FieldValue, Field.Tag)
	a.recursionLoadField(Field.Type)
}

func (a *Application) loadBeanMethod(method reflect.Value, name string) {
	methodType := method.Type()
	a.Logger.Debug("%s -> %s", name, methodType)
	if methodType.NumOut() == 1 {
		methodBean := NewBeanMethod(method, name)
		a.loadPrimeField(methodType.Out(0), name)
		for i := 0; i < methodType.NumIn(); i++ {
			inType := methodType.In(i)
			if CheckFieldPtr(inType) && ContainsFields(inType.Elem(), COMPONENT_TYPES) {
				methodBean.Ins = append(methodBean.Ins, inType)
				a.loadMethodIn(inType)
			} else {
				a.Logger.Errorf(`Param %s of %s isn't one of COMPONENT_TYPES`, inType, name)
				return
			}
		}
	}
}

// load field of configuration
func (a *Application) loadPrimeField(fieldType reflect.Type, name string) {
	if ContainsFields(fieldType.Elem(), COMPONENT_TYPES) {
		tag := (reflect.StructTag)(fmt.Sprintf(`name:"%s"`, name))
		name := GetBeanName(fieldType, tag)
		if ConfirmAddBeanMap(a.BeanMap, fieldType, name) {
			a.load(fieldType, reflect.New(fieldType.Elem()).Elem(), tag)
			a.recursionLoadField(fieldType)
		} else {
			a.Logger.Errorf("There too many %s named %s in Application", fieldType, name)
		}
	}
}

func (a *Application) loadMethodIn(inType reflect.Type) {
	a.load(inType, reflect.New(inType.Elem()).Elem(), *new(reflect.StructTag))
	a.recursionLoadField(inType)
}

func (a *Application) recursionLoadField(fieldType reflect.Type) {
	if CheckFieldPtr(fieldType) {
		appType := fieldType.Elem()
		if appType.Kind() == reflect.Struct {
			for i := 0; i < appType.NumField(); i++ {
				field := appType.Field(i)
				if CheckComponents(field) {
					a.loadField(field, reflect.New(field.Type.Elem()).Elem())
				}
			}
		}
	}
}

func (a *Application) loadConfigFile() *Application {
	appType := reflect.TypeOf(a.app).Elem()
	for i := 0; i < appType.NumField(); i++ {
		field := appType.Field(i)
		if field.Type == reflect.TypeOf(CONF{}) {
			path := DEFAULT_PATH
			tag := field.Tag
			truePath := strings.Trim(tag.Get("path"), " ")
			fileType := strings.Trim(tag.Get("type"), " ")
			if truePath != "" {
				path = truePath
			} else {
				a.Logger.Info("Conf file path unexpected, use %s", DEFAULT_PATH)
			}

			if fileType != "" && fileType != JSON && fileType != YAML {
				a.Logger.Info("Conf file suffix .%s unexpected, use default", fileType)
			}

			a.Logger.Debug("Load %s(%s)", path, fileType)

			config, err := GetJSONFromAnyFile(path, fileType)
			if err == nil {
				a.ConfigFile = config
			} else {
				a.Logger.Error("File %s load failed, %s", path, err.Error())
			}
		}
	}
	return a
}

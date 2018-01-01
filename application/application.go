package application

import (
	"reflect"
	"summer/bean"
	"github.com/labstack/echo"
	"summer/types"
	"summer"
	"summer/logger"
	"fmt"
)

type Application struct {
	app             interface{}
	BeanMap         map[reflect.Type]map[string]*bean.Bean
	BeanMethodMap   map[reflect.Type]map[string]*bean.Method
	ControllerSlice []*bean.Controller
	MiddlewareSlice []*bean.Middleware
	Server          *echo.Echo
	Logger          *logger.Logger
}

func CreateApplication(app interface{}) *Application {
	if summer.CheckFieldPtr(reflect.TypeOf(app)) {
		return (&Application{
			app:             app,
			BeanMap:         make(map[reflect.Type]map[string]*bean.Bean),
			BeanMethodMap:   make(map[reflect.Type]map[string]*bean.Method),
			ControllerSlice: make([]*bean.Controller, 0),
			MiddlewareSlice: make([]*bean.Middleware, 0),
			Server:          echo.New(),
			Logger:          logger.NewLogger(),
		}).init()
	}
	logger.NewLogger().Errorf("%s is not kind of Ptr!!!\n", reflect.TypeOf(app).Name)
	return new(Application)
}

func (a *Application) init() *Application {
	return a.loadElem(a.Logger, *new(reflect.StructTag)).loadElem(a, *new(reflect.StructTag))
}

func (a *Application) loadElem(elem interface{}, tag reflect.StructTag) *Application {
	Value := reflect.ValueOf(elem)
	Type := reflect.TypeOf(elem)
	return a.load(Type, Value.Elem(), tag)
}

func (a *Application) load(fieldType reflect.Type, Value reflect.Value, tag reflect.StructTag) *Application {
	a.Logger.Debug("%s", Value.Type())
	name := summer.GetBeanName(fieldType, tag)
	if summer.ConfirmAddBeanMap(a.BeanMap, fieldType, name) {
		newBean := bean.NewBean(Value, tag)
		a.BeanMap[fieldType][name] = newBean
	}
	return a
}

func (a *Application) Run() {
	app := a.app
	appType := reflect.TypeOf(app).Elem()
	for i := 0; i < appType.NumField(); i++ {
		field := appType.Field(i)
		if summer.CheckConfiguration(field) {
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
	a.Logger.Info("%s -> %s", name, methodType)
	if methodType.NumOut() == 1 {
		a.loadMethodOut(methodType.Out(0), name)
	}
}

func (a *Application) loadMethodOut(Out reflect.Type, name string) {
	if summer.ContainsFields(Out.Elem(), types.COMPONENT_TYPES) {
		name := summer.GetBeanName(Out, *new(reflect.StructTag))
		if a.BeanMap[Out] == nil {
			a.BeanMap[Out] = make(map[string]*bean.Bean)
		} else if _, ok := a.BeanMap[Out][name]; ok {
			a.Logger.Errorf("There too many %s named %s in Application", Out, name)
			return
		}
		a.load(Out, reflect.New(Out.Elem()).Elem(), (reflect.StructTag)(fmt.Sprintf(`name:"%s"`, name)))
		a.recursionLoadField(Out)
	}
}

func (a *Application) recursionLoadField(fieldType reflect.Type) {
	if summer.CheckFieldPtr(fieldType) {
		appType := fieldType.Elem()
		if appType.Kind() == reflect.Struct {
			for i := 0; i < appType.NumField(); i++ {
				field := appType.Field(i)
				if summer.CheckComponents(field) {
					a.loadField(field, reflect.New(field.Type.Elem()).Elem())
				}
			}
		}
	}
}

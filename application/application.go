package application

import (
	"reflect"
	"rhapsody/bean"
	"github.com/labstack/echo"
	"rhapsody/types"
	"rhapsody"
	"rhapsody/logger"
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
	if rhapsody.CheckFieldPtr(reflect.TypeOf(app)) {
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
	name := rhapsody.GetBeanName(fieldType, tag)
	if rhapsody.ConfirmAddBeanMap(a.BeanMap, fieldType, name) {
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
		if rhapsody.CheckConfiguration(field) {
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
		methodBean := bean.NewBeanMethod(method, name)
		a.loadMethodOut(methodType.Out(0), name)
		for i:=0; i < methodType.NumIn(); i++ {
			inType := methodType.In(i)
			if rhapsody.CheckFieldPtr(inType) && rhapsody.ContainsFields(inType.Elem(), types.COMPONENT_TYPES) {
				methodBean.Ins = append(methodBean.Ins, inType)
				a.loadMethodIn(inType)
			} else {
				a.Logger.Errorf(`Param %s of %s isn't one of COMPONENT_TYPES`, inType, name)
				return
			}
		}
	}
}

func (a *Application) loadMethodOut(Out reflect.Type, name string) {
	if rhapsody.ContainsFields(Out.Elem(), types.COMPONENT_TYPES) {
		tag := (reflect.StructTag)(fmt.Sprintf(`name:"%s"`, name))
		name := rhapsody.GetBeanName(Out, tag)
		if a.BeanMap[Out] == nil {
			a.BeanMap[Out] = make(map[string]*bean.Bean)
		} else if _, ok := a.BeanMap[Out][name]; ok {
			a.Logger.Errorf("There too many %s named %s in Application", Out, name)
			return
		}
		a.load(Out, reflect.New(Out.Elem()).Elem(), tag)
		a.recursionLoadField(Out)
	}
}

func (a *Application) loadMethodIn(inType reflect.Type) {
	a.load(inType, reflect.New(inType.Elem()).Elem(), *new(reflect.StructTag))
	a.recursionLoadField(inType)
}

func (a *Application) recursionLoadField(fieldType reflect.Type) {
	if rhapsody.CheckFieldPtr(fieldType) {
		appType := fieldType.Elem()
		if appType.Kind() == reflect.Struct {
			for i := 0; i < appType.NumField(); i++ {
				field := appType.Field(i)
				if rhapsody.CheckComponents(field) {
					a.loadField(field, reflect.New(field.Type.Elem()).Elem())
				}
			}
		}
	}
}

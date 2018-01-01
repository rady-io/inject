package application

import (
	"reflect"
	"summer/bean"
	"github.com/labstack/echo"
	"summer/types"
	"summer"
)

type Application struct {
	app           interface{}
	BeanMap       map[reflect.Type]map[string]*bean.Bean
	BeanMethodMap map[reflect.Type]map[string]*bean.Method
	ControllerSlice []*bean.Controller
	MiddlewareSlice []*bean.Middleware
	Server        *echo.Echo
}

func CreateApplication(app interface{}) *Application {
	return &Application{
		app:           app,
		BeanMap:       make(map[reflect.Type]map[string]*bean.Bean),
		BeanMethodMap: make(map[reflect.Type]map[string]*bean.Method),
		ControllerSlice: make([]*bean.Controller, 0),
		MiddlewareSlice: make([]*bean.Middleware, 0),
		Server:        echo.New(),
	}
}

func (a *Application) Run() {
	app := a.app
	appValue := reflect.Indirect(reflect.ValueOf(app)) // can get Elem if app is pointer of a struct
	appType := appValue.Type()
	for i := 0; i < appType.NumField(); i++ {
		field := appType.Field(i)
		fieldValue := appValue.Field(i)
		if field.Tag.Get("type") == "" && summer.ContainsField(reflect.Indirect(reflect.New(field.Type)).Type(), types.Configuration{}) || field.Tag.Get("type") == types.CONFIGURATION {
			a.loadField(field, fieldValue)
		}
	}
}

func (a *Application) loadField(Field reflect.StructField, FieldValue reflect.Value) {
	fieldType := Field.Type
	if a.BeanMap[fieldType] == nil {
		newBean := bean.NewBean(FieldValue, Field.Tag)
		a.BeanMap[fieldType] = make(map[string]*bean.Bean)
		a.BeanMap[fieldType][fieldType.Name()] = newBean
		a.recursion(FieldValue)
	}
	FieldValue.Set(a.BeanMap[fieldType][fieldType.Name()].Value)
}

func (a *Application) recursion(value reflect.Value) {
	appValue := reflect.Indirect(value) // can get Elem if app is pointer of a struct
	appType := appValue.Type()
	for i := 0; i < appType.NumField(); i++ {
		field := appType.Field(i)
		if _, ok := types.COMPONENTS[field.Tag.Get("type")]; field.Tag.Get("type") == "" && summer.ContainsFields(field.Type, types.COMPONENT_TYPES) || ok {
			a.loadField(field, appValue.Field(i))
		}
	}
}

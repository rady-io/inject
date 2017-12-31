package application

import (
	"reflect"
	"summer/bean"
	"github.com/labstack/echo"
	"summer/types"
	"summer"
)

type Application struct {
	app       interface{}
	beamMap   map[reflect.Type]map[string]*bean.Bean
	methodMap map[reflect.Type]map[string]*bean.Method
	Server    *echo.Echo
}

func CreateApplication(app interface{}) *Application {
	return &Application{
		app:       app,
		beamMap:   make(map[reflect.Type]map[string]*bean.Bean),
		methodMap: make(map[reflect.Type]map[string]*bean.Method),
		Server:    echo.New(),
	}
}

func (a *Application) Run() {
	app := a.app
	appValue := reflect.Indirect(reflect.ValueOf(app)) // can get Elem if app is pointer of a struct
	appType := appValue.Type()
	for i := 0; i < appType.NumField(); i++ {
		field := appType.Field(i)
		if field.Tag.Get("type") == "" && summer.ContainsField(field.Type, types.Configuration{}) || field.Tag.Get("type") == types.CONFIGURATION {
			a.loadField(field, appValue.Field(i))
		}
	}
}

func (a *Application) loadField(Field reflect.StructField, FieldValue reflect.Value) {
	fieldType := Field.Type
	if a.beamMap[fieldType] == nil {
		newBean := bean.NewBean(FieldValue, Field.Tag)
		a.beamMap[fieldType] = make(map[string]*bean.Bean)
		a.beamMap[fieldType][fieldType.Name()] = newBean
	}
	FieldValue.Set(a.beamMap[fieldType][fieldType.Name()].Value)
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

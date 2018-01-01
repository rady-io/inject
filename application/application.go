package application

import (
	"reflect"
	"summer/bean"
	"github.com/labstack/echo"
	"summer/types"
	"summer"
	"summer/logger"
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
	return a.initialize(a.Logger).initialize(a)
}

func (a *Application) initialize(elem interface{}) *Application {
	Value := reflect.ValueOf(elem)
	Type := reflect.TypeOf(elem)
	a.BeanMap[Type] = map[string]*bean.Bean{
		Type.Name(): bean.NewBean(Value, *new(reflect.StructTag)),
	}
	return a
}

func (a *Application) Run() {
	app := a.app
	appValue := reflect.ValueOf(app).Elem() // can get Elem if app is pointer of a struct
	appType := reflect.TypeOf(app).Elem()
	for i := 0; i < appType.NumField(); i++ {
		field := appType.Field(i)
		fieldValue := appValue.Field(i)
		if summer.CheckFieldPtr(field.Type) && (field.Tag.Get("type") == "" && summer.ContainsField(field.Type.Elem(), types.Configuration{}) || field.Tag.Get("type") == types.CONFIGURATION) {
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
		a.recursionLoadField(Field)
	}
	//FieldValue.Set(a.BeanMap[fieldType][fieldType.Name()].Value)
}

func (a *Application) recursionLoadField(Field reflect.StructField) {
	if summer.CheckFieldPtr(Field.Type) {
		//appValue := FieldValue.Elem() // can get Elem if app is pointer of a struct
		appType := Field.Type.Elem()
		for i := 0; i < appType.NumField(); i++ {
			field := appType.Field(i)
			_, ok := types.COMPONENTS[field.Tag.Get("type")]
			if summer.CheckFieldPtr(field.Type) && (ok || field.Tag.Get("type") == "" && summer.ContainsFields(field.Type.Elem(), types.COMPONENT_TYPES)) {
				a.loadField(field, reflect.New(field.Type))
			}
		}
	}
}

package application

import (
	"reflect"
	"summer/bean"
	"github.com/labstack/echo"
	"summer/types"
	"summer"
	"summer/logger"
	"strings"
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
	return a.load(Type, Value, tag)
}

func (a *Application) load(Type reflect.Type, Value reflect.Value, tag reflect.StructTag) *Application {
	name := Type.String()
	if tag != *new(reflect.StructTag) {
		if aliasName := tag.Get("name"); strings.Trim(aliasName, " ") != "" {
			name = aliasName
		}
	}
	a.BeanMap[Type] = map[string]*bean.Bean{
		name: bean.NewBean(Value, tag),
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
		a.load(fieldType, FieldValue, Field.Tag)
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

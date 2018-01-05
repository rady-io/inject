package rady

import (
	"reflect"
	"github.com/labstack/echo"
	"strings"
	"os"
)

/*
Application is the bootstrap of a Rhapsody app

Root is pointer of app for config, controller and handler registry

BeanMap is map to find *Bean by `Type` and `Name`(when type is the same)

	value in bean is Elem(), can get Addr() when set to other field

BeanMethodMap is map to find *Method by `Type` and `Name`(when type is the same)

	Value in method can be call later to set value to other filed

ValueBeanMap is map to find / set value in config file, going to implement hot-reload

CtrlBeanSlice is slice to store controller

MdWareBeanSlice is slice to store middleware

Entities is slice to store type of entity

Server is the echo server

Logger is the global logger

ConfigFile is the string json value of config file
*/
type Application struct {
	Root            interface{}
	BeanMap         map[reflect.Type]map[string]*Bean
	BeanMethodMap   map[reflect.Type]map[string]*Method
	ValueBeanMap    map[string]*ValueBean
	CtrlBeanSlice   []*CtrlBean
	MdWareBeanSlice []*MdWareBean
	Entities        []reflect.Type
	Server          *echo.Echo
	Logger          *Logger
	ConfigFile      string
}

/*
CreateApplication can initial application with root

if root is not kinds of Ptr, there will be an error
 */
func CreateApplication(root interface{}) *Application {
	if CheckFieldPtr(reflect.TypeOf(root)) {
		return (&Application{
			Root:            root,
			BeanMap:         make(map[reflect.Type]map[string]*Bean),
			BeanMethodMap:   make(map[reflect.Type]map[string]*Method),
			ValueBeanMap:    make(map[string]*ValueBean),
			CtrlBeanSlice:   make([]*CtrlBean, 0),
			MdWareBeanSlice: make([]*MdWareBean, 0),
			Entities:        make([]reflect.Type, 0),
			Server:          echo.New(),
			Logger:          NewLogger(),
		}).init()
	}
	NewLogger().Errorf("%s is not kind of Ptr!!!\n", reflect.TypeOf(root).Name())
	return new(Application)
}

func (a *Application) init() *Application {
	return a.loadElem(a.Logger, *new(reflect.StructTag)).loadElem(a, *new(reflect.StructTag)).loadConfigFile()
}

func (a *Application) loadElem(elem interface{}, tag reflect.StructTag) *Application {
	Value := reflect.ValueOf(elem)
	Type := reflect.TypeOf(elem)
	return a.LoadBean(Type, Value.Elem(), tag)
}

func (a *Application) load(fieldType reflect.Type, Value reflect.Value, tag reflect.StructTag) *Application {
	name := GetBeanName(fieldType, tag)
	a.Logger.Debug("%s -> %s", name, Value.Type())
	newBean := NewBean(Value, tag)
	a.BeanMap[fieldType][name] = newBean
	return a
}

/*
Run is the boot method of a whole Rhapsody app

Start with the Root, Get Fields of Root

First, we load "PrimeBean", which mean field define in config, are basic of all bean

	PrimeBean include all component Fields in config, which can defined unique name

Then, we load BeanMethodOut, which mean result of "bean method", however, what's bean method?

	Bean method mean methods defined in config, return only bean(component), and all parameters can init with dependency injection

	returned value of Bean method is same as PrimeBean

	parameter value is same as normal field

And then, we load normal bean recursively

TODO: Inject Value

TODO: Load BeanMethod Args, Construct Link between Param and Value

TODO: Load Controller and Middleware and then register them

 */
func (a *Application) Run() {
	root := a.Root
	rootType := reflect.TypeOf(root).Elem()
	for i := 0; i < rootType.NumField(); i++ {
		config := rootType.Field(i)
		if CheckConfiguration(config) {
			configValue := reflect.New(config.Type.Elem()).Elem() // save Elem in Bean
			a.LoadBean(config.Type, configValue, config.Tag)
			for i := 0; i < configValue.NumField(); i++ {
				fieldValue := configValue.Field(i)
				field := config.Type.Elem().Field(i)
				if CheckConfiguration(field) {
					a.LoadPrimeBean(field.Type, fieldValue, field.Tag)
				}
			}
			for i := 0; i < configValue.Addr().NumMethod(); i++ {
				name := configValue.Addr().Type().Method(i).Name
				a.loadBeanMethodOut(configValue.Addr().MethodByName(name), name)
			}
		}
	}
	a.loadMethodBeanIn()
	a.loadBeanChild()
}

// load children of a Bean
func (a *Application) loadBeanChild() {
	for fileType, beanMap := range a.BeanMap {
		if len(beanMap) > 0 {
			a.RecursivelyLoad(fileType)
		}
	}
}

func (a *Application) loadMethodBeanIn() {
	for _, methodMap := range a.BeanMethodMap {
		for _, method := range methodMap {
			method.LoadIns(a)
		}
	}
}

/*
LoadBean can load normal bean

normal bean can have a name, only if there is a prime bean with same name and type, and this method will do nothing

if a normal bean doesn't a name

	1. there is only one loaded bean with the same type, this method do nothing

	2. there are more than one loaded bean with the same type, error and exit

	3. there is no loaded bean with the same type, this method with initialize a new bean
 */
func (a *Application) LoadBean(fieldType reflect.Type, fieldValue reflect.Value, tag reflect.StructTag) *Application {
	name := tag.Get("name")
	if name == "" {
		if ConfirmSameTypeInMap(a.BeanMap, fieldType) {
			if len(a.BeanMap[fieldType]) > 1 {
				a.Logger.Critical("There are more than one %s, please named it.", fieldType)
				os.Exit(1)
			}
		} else {
			a.load(fieldType, fieldValue, tag)
		}
	} else {
		if !ConfirmBeanInMap(a.BeanMap, fieldType, name) {
			a.Logger.Critical("There are no %s named %s, please define it in config.", fieldType, name)
			os.Exit(1)
		}
	}
	return a
}

func (a *Application) loadBeanMethodOut(method reflect.Value, name string) {
	methodType := method.Type()
	a.Logger.Debug("%s -> %s", name, methodType)
	if methodType.NumOut() == 1 {
		methodBean := NewBeanMethod(method, name)
		fieldType := methodType.Out(0)
		OutValue := reflect.New(fieldType.Elem()).Elem()
		if a.LoadPrimeBean(fieldType, OutValue, GetTagFromName(name)) {
			for i := 0; i < methodType.NumIn(); i++ {
				inType := methodType.In(i)
				if CheckFieldPtr(inType) && ContainsFields(inType.Elem(), ComponentTypes) {
					methodBean.Ins = append(methodBean.Ins, inType)
				} else {
					a.Logger.Errorf(`Param %s of %s isn't one of ComponentTypes`, inType, name)
					os.Exit(1)
				}
			}
			methodBean.OutValue = OutValue
			a.BeanMethodMap[fieldType] = map[string]*Method{name: methodBean}
		}
	}
}

/*
LoadPrimeBean as its name

load field in configuration and out of beanMethod in configuration
 */
func (a *Application) LoadPrimeBean(fieldType reflect.Type, fieldValue reflect.Value, tag reflect.StructTag) bool {
	if ContainsFields(fieldType.Elem(), ComponentTypes) {
		tag := tag
		name := GetBeanName(fieldType, tag)
		if ConfirmAddBeanMap(a.BeanMap, fieldType, name) {
			a.load(fieldType, fieldValue, tag)
			return true
		} else {
			a.Logger.Errorf("There too many %s named %s in Application", fieldType, name)
		}
	}
	return false
}

/*
LoadBeanAndRecursion initialize a instance of a type and load it

then recursively load children of this type
 */
func (a *Application) LoadBeanAndRecursion(fieldType reflect.Type) {
	a.LoadBean(fieldType, reflect.New(fieldType.Elem()).Elem(), *new(reflect.StructTag))
	a.RecursivelyLoad(fieldType)
}

/*
RecursivelyLoad recursively load children of a normal bean

only if there is no prime bean among the children
 */
func (a *Application) RecursivelyLoad(fieldType reflect.Type) {
	if CheckFieldPtr(fieldType) && ContainsFields(fieldType.Elem(), ComponentTypes) {
		appType := fieldType.Elem()
		if appType.Kind() == reflect.Struct {
			for i := 0; i < appType.NumField(); i++ {
				field := appType.Field(i)
				if CheckComponents(field) {
					a.LoadBean(field.Type, reflect.New(field.Type.Elem()).Elem(), field.Tag)
				}
			}
		}
	}
}

func (a *Application) loadConfigFile() *Application {
	appType := reflect.TypeOf(a.Root).Elem()
	for i := 0; i < appType.NumField(); i++ {
		field := appType.Field(i)
		if field.Type == reflect.TypeOf(CONF{}) {
			path := DefaultPath
			tag := field.Tag
			truePath := strings.Trim(tag.Get("path"), " ")
			fileType := strings.Trim(tag.Get("type"), " ")
			if truePath != "" {
				path = truePath
			} else {
				a.Logger.Info("Conf file path unexpected, use %s", DefaultPath)
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

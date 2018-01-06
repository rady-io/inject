package rady

import (
	"reflect"
	"github.com/labstack/echo"
	"strings"
	"os"
	"github.com/tidwall/gjson"
	"fmt"
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

 */
func (a *Application) Run() {
	a.loadPrimes()
	a.loadMethodBeanIn()
	a.loadBeanChild()
	a.assemble()
	a.CallFactory()
	a.Server.Start(a.getAddr())
}

func (a *Application) getAddr() string {
	result := gjson.Get(a.ConfigFile, "rhapsody.server.addr")
	if result.Exists() {
		return result.String()
	}
	return ":8081"
}

func (a *Application) loadPrimes() {
	root := a.Root
	rootType := reflect.TypeOf(root).Elem()
	for i := 0; i < rootType.NumField(); i++ {
		field := rootType.Field(i)
		if CheckConfiguration(field) {
			a.loadConfiguration(field)
		} else {
			a.loadWebField(field, "/")
		}
	}
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

func (a *Application) loadWebField(field reflect.StructField, prefix string) {
	if CheckRouter(field) {
		a.loadRouter(field, prefix)
	} else if CheckController(field) {
		a.loadCtrl(field, prefix)
	} else if CheckMiddleware(field) {
		a.loadMiddleware(field, prefix)
	}
}

func (a *Application) loadRouter(field reflect.StructField, prefix string) {
	path := GetPathFromType(field, Router{})
	prefix = GetNewPrefix(prefix, path)
	fieldType := field.Type

	for i := 0; i < fieldType.Elem().NumField(); i++ {
		a.loadWebField(fieldType.Elem().Field(i), prefix)
	}
}

func (a *Application) loadCtrl(field reflect.StructField, prefix string) {
	loadedMethod := make(map[string]bool)
	path := GetPathFromType(field, Controller{})
	prefix = GetNewPrefix(prefix, path)
	fieldType := field.Type
	Name := field.Name
	Value := reflect.New(fieldType.Elem()).Elem()
	a.CtrlBeanSlice = append(a.CtrlBeanSlice, NewCtrlBean(Value, field.Tag, Name))
	a.LoadPrimeBean(fieldType, Value, ``)

	for i := 0; i < fieldType.Elem().NumField(); i++ {
		child := fieldType.Elem().Field(i)
		if httpMethod, ok := MethodsTypeSet[child.Type]; ok {
			path := GetNewPrefix(prefix, child.Tag.Get("path"))
			handlerName := child.Tag.Get("method")
			if _, ok := loadedMethod[handlerName]; !ok && handlerName != "" {
				if method := Value.Addr().MethodByName(handlerName); method.IsValid() {
					if trueMethod, ok := method.Interface().(func(Context) error); ok {
						a.registerCtrl(trueMethod, child.Type, path)
						a.logHandlerRegistry(httpMethod, path, handlerName)
						loadedMethod[handlerName] = true
					}
				}
			}
		} else if child.Type == StaticType {
			a.loadStatic(child, prefix)
		} else if child.Type == FileType {
			a.loadFile(child, prefix)
		}
	}

	for i := 0; i < fieldType.NumMethod(); i++ {
		method := Value.Addr().Method(i)
		methodField := fieldType.Method(i)
		handlerName := methodField.Name
		if _, ok := loadedMethod[handlerName]; !ok {
			ok, httpMethod, path := ParseHandlerName(handlerName)
			if ok {
				if trueMethod, ok := method.Interface().(func(Context) error); ok {
					path := GetNewPrefix(prefix, path)
					a.registerCtrl(trueMethod, reflect.TypeOf(httpMethod), path)
					a.logHandlerRegistry(MethodToStr[httpMethod], path, handlerName)
					loadedMethod[handlerName] = true
				}
			}
		}
	}
}

func (a *Application) loadMiddleware(field reflect.StructField, prefix string) {
	path := GetPathFromType(field, Middleware{})
	prefix = GetNewPrefix(prefix, path)
	fieldType := field.Type
	Name := field.Name
	Value := reflect.New(fieldType.Elem()).Elem()
	a.MdWareBeanSlice = append(a.MdWareBeanSlice, NewMdWareBean(Value, field.Tag, Name))
	a.LoadPrimeBean(fieldType, Value, ``)

	for i := 0; i < fieldType.NumMethod(); i++ {
		method := Value.Addr().Method(i)
		methodField := fieldType.Method(i)
		handlerName := methodField.Name
		if trueMethod, ok := method.Interface().(func(handlerFunc HandlerFunc) HandlerFunc); ok {
			a.Server.Group(prefix, trueMethod)
			a.logMiddlewareRegistry(prefix, handlerName)
		}
	}
}

func (a *Application) loadStatic(field reflect.StructField, prefix string) {
	tag := field.Tag
	prefix = GetNewPrefix(prefix, tag.Get("prefix"))
	root := tag.Get("root")
	a.Server.Static(prefix, root)
	a.Logger.Debug("Register Static >>> %s <- %s", prefix, root)
}

func (a *Application) loadFile(field reflect.StructField, prefix string) {
	tag := field.Tag
	prefix = GetNewPrefix(prefix, tag.Get("path"))
	file := strings.Trim(tag.Get("file"), " ")
	if !CheckFilenameValid(file) {
		a.Logger.Error("File name '%s' invalid", file)
		return
	}
	a.Server.File(prefix, file)
	a.Logger.Debug("Register Static >>> %s <- %s", prefix, file)
}

func (a *Application) loadConfiguration(config reflect.StructField) {
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
				return a
			} else {
				a.Logger.Error("File %s load failed, %s", path, err.Error())
			}
		}
	}

	config, err := GetJSONFromAnyFile(DefaultPath, JSON)
	if err == nil {
		a.Logger.Debug("Load %s(%s)", DefaultPath, JSON)
		a.ConfigFile = config
	}
	return a
}

func (a *Application) logHandlerRegistry(method string, path, Name string) {
	a.Logger.Debug("Register Handler: %s >>> %s %s", Name, method, path)
}

func (a *Application) logMiddlewareRegistry(path, Name string) {
	a.Logger.Debug("Register Middleware: %s >>> %s", Name, path)
}

func (a *Application) registerCtrl(handlerFunc HandlerFunc, method reflect.Type, path string) {
	switch method {
	case reflect.TypeOf(GET{}):
		a.Server.GET(path, handlerFunc)
	case reflect.TypeOf(POST{}):
		a.Server.POST(path, handlerFunc)
	case reflect.TypeOf(PUT{}):
		a.Server.PUT(path, handlerFunc)
	case reflect.TypeOf(DELETE{}):
		a.Server.DELETE(path, handlerFunc)
	case reflect.TypeOf(OPTIONS{}):
		a.Server.OPTIONS(path, handlerFunc)
	case reflect.TypeOf(GET{}):
		a.Server.GET(path, handlerFunc)
	case reflect.TypeOf(HEAD{}):
		a.Server.HEAD(path, handlerFunc)
	case reflect.TypeOf(CONNECT{}):
		a.Server.CONNECT(path, handlerFunc)
	case reflect.TypeOf(TRACE{}):
		a.Server.TRACE(path, handlerFunc)
	case reflect.TypeOf(OPTIONS{}):
		a.Server.OPTIONS(path, handlerFunc)
	case reflect.TypeOf(PATCH{}):
		a.Server.PATCH(path, handlerFunc)
	default:
		break
	}
}

func (a *Application) assemble() {
	for beanType, nameMap := range a.BeanMap {
		for name, bean := range nameMap {
			Value := bean.Value
			for i := 0; i < beanType.Elem().NumField(); i++ {
				child := beanType.Elem().Field(i)
				if CheckComponents(child) {
					a.assembleBean(name, beanType.String(), Value.Field(i), child)
				} else if CheckValues(child) {
					a.assembleValue(name, beanType.String(), Value.Field(i), child)
				}
			}
		}
	}
}

func (a *Application) logAssembleBean(motherName, motherType, childName, childType, fieldName string) {
	a.Logger.Debug("Field %s of [%s][%s] set [%s][%s]", fieldName, motherType, motherName, childType, childName)
}

func (a *Application) logAssembleValue(motherName, motherType, value, fieldName string) {
	a.Logger.Debug("Field %s of [%s][%s] set value &(%s)", fieldName, motherType, motherName, value)
}

func (a *Application) assembleBean(motherName, motherType string, value reflect.Value, field reflect.StructField) {
	if len(a.BeanMap[field.Type]) == 0 {
		a.Logger.Critical("Type %s doesn't exist in beanMap !!!", field.Type)
		os.Exit(1)
	}

	if len(a.BeanMap[field.Type]) == 1 {
		for name, bean := range a.BeanMap[field.Type] {
			value.Set(bean.Value.Addr())
			a.logAssembleBean(motherName, motherType, name, field.Type.String(), field.Name)
		}
	}

	if len(a.BeanMap[field.Type]) > 1 {
		name := GetBeanName(field.Type, field.Tag)
		if bean, ok := a.BeanMap[field.Type][name]; ok {
			value.Set(bean.Value.Addr())
			a.logAssembleBean(motherName, motherType, name, field.Type.String(), field.Name)
		} else {
			a.Logger.Critical("Bean [%s][%s] doesn't exist in beanMap !!!", field.Type, name)
			os.Exit(1)
		}
	}
}

func (a *Application) assembleValue(motherName, motherType string, value reflect.Value, field reflect.StructField) {
	key := strings.Trim(field.Tag.Get("value"), " ")
	defaultValue := field.Tag.Get("default")
	if valueBean, ok := a.ValueBeanMap[key]; ok {
		if !valueBean.SetValue(value, field.Type) {
			a.Logger.Error("Unknow Type: %s", field.Type)
		}
		a.logAssembleValue(motherName, motherType, valueBean.Value.String(), field.Name)
		return
	}

	newValue := gjson.Get(a.ConfigFile, key)
	if !newValue.Exists() {
		newValue = gjson.Get(fmt.Sprintf(`{"default": "%s"}`, defaultValue), "default")
	}

	valueBean := NewValueBean(newValue, key)
	if !valueBean.SetValue(value, field.Type) {
		a.Logger.Error("Unknow Type: %s", field.Type)
	}

	a.ValueBeanMap[key] = valueBean
	a.logAssembleValue(motherName, motherType, valueBean.Value.String(), field.Name)
}

func (a *Application) CallFactory() {
	for outType, methodMap := range a.BeanMethodMap {
		for name, method := range methodMap {
			method.Call(a)
			a.Logger.Debug("Call factory %s -> %s", name, outType)
		}
	}
}

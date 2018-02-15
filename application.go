package rady

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/tidwall/gjson"
	"os"
	"reflect"
	"strings"
	"testing"
)

/*
Application is the bootstrap of a Rady app

Root is pointer of app for config, controller and handler registry

BeanMap is map to find *Bean by `Type` and `Name`(when type is the same)

	value in bean is Elem(), can get Addr() when set to other field

BeanMethodMap is map to find *Method by `Type` and `Name`(when type is the same)

	Value in method can be call later to set value to other filed

ValueBeanMap is map to find / set value in config file, going to implement hot-reload

CtrlBeanMap is slice to store controller

MdWareBeanMap is slice to store middleware

Entities is slice to store type of entity

Server is the echo server

Logger is the global logger

ConfigFile is the string json value of config file
*/
type Application struct {
	BootStrap
	Root               interface{}
	BeanMap            map[reflect.Type]map[string]*Bean
	BeanMethodMap      map[reflect.Type]map[string]*Method
	ValueBeanMap       map[string]*ValueBean
	FactoryToRecall    map[*Method]bool
	CtrlBeanMap        map[string]*CtrlBean
	MdWareBeanMap      map[string]*MdWareBean
	MiddlewareStackMap map[string]*MiddlewareStack
	Entities           []reflect.Type
	TestingBeans       []*TestingBean
	Server             *echo.Echo
	Logger             *Logger
	ConfigFile         string
	Addr               *string `value:"rady.server.addr" default:":8081"`
}

/*
CreateApplication can initial application with root

if root is not kinds of Ptr, there will be an error
*/
func CreateApplication(root interface{}) *Application {
	if CheckPtrOfStruct(reflect.TypeOf(root)) {
		return (&Application{
			Root:               root,
			BeanMap:            make(map[reflect.Type]map[string]*Bean),
			BeanMethodMap:      make(map[reflect.Type]map[string]*Method),
			ValueBeanMap:       make(map[string]*ValueBean),
			FactoryToRecall:    make(map[*Method]bool),
			CtrlBeanMap:        make(map[string]*CtrlBean),
			MdWareBeanMap:      make(map[string]*MdWareBean),
			MiddlewareStackMap: make(map[string]*MiddlewareStack),
			Entities:           make([]reflect.Type, 0),
			TestingBeans:       make([]*TestingBean, 0),
			Server:             echo.New(),
			Logger:             NewLogger(),
		}).init()
	}
	NewLogger().Errorf("%s is not kind of Ptr!!!\n", reflect.TypeOf(root).Name())
	return new(Application)
}

func CreateTest(root interface{}) *Application {
	return CreateApplication(root).PrepareTest()
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

func (a *Application) Prepare() *Application {
	a.loadPrimes()
	a.loadMethodBeanIn()
	return a
}

func (a *Application) PrepareTest() *Application {
	ComponentTypes[reflect.TypeOf(Testing{})] = true
	COMPONENTS[TESTING] = true
	return a.Prepare()
}

func (a *Application) setTests(testType reflect.Type, testValue reflect.Value) {
	for i := 0; i < testValue.NumField(); i++ {
		field := testType.Elem().Field(i)
		if CheckTesting(field) {
			fieldValue := reflect.New(field.Type.Elem()).Elem()
			a.setTest(field.Type, fieldValue, field.Tag)
		}
	}
}

func (a *Application) setTest(testType reflect.Type, testValue reflect.Value, Tag reflect.StructTag) {
	a.LoadBean(testType, testValue, Tag)
	a.TestingBeans = append(a.TestingBeans, NewTestingBean(testType, testValue))
	a.Logger.Debug("SetTest: %s", testType)
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

*/
func (a *Application) Run() {
	a.Prepare()
	a.loadBeanChild()
	a.assemble()
	a.CallFactory()
	a.bindFactoryWithValue()
	a.Server.Start(*a.Addr)
}

func (a *Application) Test(t *testing.T) *Application {
	a.loadBeanChild()
	a.assemble()
	a.CallFactory()
	a.bindFactoryWithValue()
	a.runTestCase(t)
	return a
}

func (a *Application) AddTest(testPointer interface{}) *Application {
	testType := reflect.TypeOf(testPointer)
	if !CheckPtrOfStruct(testType) {
		return a
	}
	testValue := reflect.New(testType.Elem()).Elem()
	a.setTest(testType, testValue, "")
	return a
}

func (a *Application) AddTests(testPointer interface{}) *Application {
	testType := reflect.TypeOf(testPointer)
	if !CheckPtrOfStruct(testType) {
		return a
	}
	testValue := reflect.New(testType.Elem()).Elem()
	a.setTests(testType, testValue)
	return a
}

func (a *Application) runTestCase(t *testing.T) {
	for _, testingBean := range a.TestingBeans {
		for i := 0; i < testingBean.Type.NumMethod(); i++ {
			method := testingBean.Type.Method(i)
			methodValue := testingBean.Value.Addr().Method(i)
			TestCase, ok := methodValue.Interface().(func(t *testing.T))
			if strings.HasPrefix(method.Name, "Test") && ok {
				TestCase(t)
			}
		}
	}
}

func (a *Application) loadPrimes() {
	root := a.Root
	rootType := reflect.TypeOf(root).Elem()
	for i := 0; i < rootType.NumField(); i++ {
		field := rootType.Field(i)
		if CheckConfiguration(field) {
			a.loadConfiguration(field)
		} else if CheckEntities(field) {
			a.loadEntities(field)
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
	a.CtrlBeanMap[prefix] = NewCtrlBean(Value, field.Tag, Name)
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
	newPrefix := GetNewPrefix(prefix, path)

	if a.MiddlewareStackMap[newPrefix] == nil {
		a.MiddlewareStackMap[newPrefix] = NewMiddlewareStack()
	}

	fieldType := field.Type
	Name := field.Name
	Value := reflect.New(fieldType.Elem()).Elem()
	a.MdWareBeanMap[newPrefix] = NewMdWareBean(Value, field.Tag, Name)
	a.LoadPrimeBean(fieldType, Value, ``)

	for i := 0; i < fieldType.NumMethod(); i++ {
		method := Value.Addr().Method(i)
		methodField := fieldType.Method(i)
		handlerName := methodField.Name
		if trueMethod, ok := method.Interface().(func(handlerFunc HandlerFunc) HandlerFunc); ok {
			a.MiddlewareStackMap[newPrefix] = a.MiddlewareStackMap[newPrefix].Push(NewMiddlewareContainer(handlerName, trueMethod))
			a.Server.Group(newPrefix, trueMethod)
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

	a.Logger.Debug("Load Configuration: %s", config.Type)
	for i := 0; i < configValue.Addr().NumMethod(); i++ {
		name := configValue.Addr().Type().Method(i).Name
		a.loadBeanMethodOut(configValue.Addr().MethodByName(name), name)
	}
}

func (a *Application) loadEntities(field reflect.StructField) {
	valueType := field.Type.Elem()
	for i := 0; i < valueType.NumField(); i++ {
		childType := valueType.Field(i).Type
		if childType.Kind() == reflect.Ptr && childType.Elem().Kind() == reflect.Struct {
			a.Entities = append(a.Entities, childType)
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
	if methodType.NumOut() == 1 {
		fieldType := methodType.Out(0)
		if fieldType.Kind() == reflect.Ptr && ContainsFields(fieldType.Elem(), ComponentTypes) {
			a.Logger.Debug("%s -> %s", name, methodType)
			methodBean := NewBeanMethod(method, name)
			a.Logger.Debug("Method: %s", name)
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
			os.Exit(1)
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

func (a *Application) GetRealConfigPathAndType() (string, string) {
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
				a.Logger.Info("Conf file path unexpected, use %s", truePath)
			}

			path = GetConfigFileByMode(path)

			if fileType != "" && fileType != JSON && fileType != YAML {
				a.Logger.Info("Conf file suffix .%s unexpected, use default", fileType)
				return path, JSON
			}
			return path, fileType
		}
	}

	return DefaultPath, JSON
}

func (a *Application) loadConfigFile() *Application {
	path, fileType := a.GetRealConfigPathAndType()
	config, err := GetJSONFromAnyFile(path, fileType)
	if err == nil {
		a.ConfigFile = config
		a.Logger.Debug("Load %s(%s)", path, fileType)
	} else {
		a.Logger.Error("File %s load failed, %s", path, err.Error())
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
	MethodName, ok := MethodsTypeSet[method]
	if ok {
		ServerVal := reflect.ValueOf(a.Server)
		MethodVal := ServerVal.MethodByName(strings.ToUpper(MethodName))
		MethodVal.Call([]reflect.Value{
			reflect.ValueOf(path),
			reflect.ValueOf(handlerFunc),
		})
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
	if key == "" {
		return
	}
	defaultValue := field.Tag.Get("default")
	if valueBean, ok := a.ValueBeanMap[key]; ok {
		if !valueBean.SetValue(value, field.Type) {
			a.Logger.Error("Unknow Type: %s", field.Type)
		}
		a.logAssembleValue(motherName, motherType, valueBean.Value.String(), field.Name)
		return
	}

	newValue := gjson.Get(a.ConfigFile, key)
	trueDefault := gjson.Get(fmt.Sprintf(`{"default": "%s"}`, defaultValue), "default")
	if !newValue.Exists() {
		newValue = trueDefault
	}
	valueBean := NewValueBean(newValue, key, trueDefault)
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

func (a *Application) bindFactoryWithValue() {
	checkedMap := make(map[reflect.Type]bool)
	for _, methodMap := range a.BeanMethodMap {
		for _, method := range methodMap {
			for _, inType := range method.Ins {
				a.recursivelyBind(inType, method, checkedMap)
			}
		}
	}
}

func (a *Application) recursivelyBind(fieldType reflect.Type, method *Method, checkedMap map[reflect.Type]bool) {
	checkedMap[fieldType] = true
	for i := 0; i < fieldType.Elem().NumField(); i++ {
		child := fieldType.Elem().Field(i)
		if CheckComponents(child) {
			if _, ok := checkedMap[child.Type]; !ok {
				a.recursivelyBind(child.Type, method, checkedMap)
			}
		} else if CheckPtrValues(child) {
			key := strings.Trim(child.Tag.Get("value"), " ")
			valueBean, ok := a.ValueBeanMap[key]
			if !ok {
				a.Logger.Critical("Ignored value '%s' !!!", key)
				os.Exit(1)
			}

			if _, ok = valueBean.MethodSet[method]; ok {
				continue
			}
			a.Logger.Debug("Bind value '%s' with factory %s", key, method.Name)
			valueBean.MethodSet[method] = true
		}
	}
}

func (a *Application) ReloadValues() {
	a.loadConfigFile()
	a.FactoryToRecall = make(map[*Method]bool)
	for _, valueBean := range a.ValueBeanMap {
		valueBean.Reload(a)
	}

	for recallFactory := range a.FactoryToRecall {
		recallFactory.Call(a)
	}
}

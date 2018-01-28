package rady

import "reflect"

func init() {
	COMPONENTS[COMPONENT] = true
	COMPONENTS[SERVICE] = true
	COMPONENTS[REPOSITORY] = true
	COMPONENTS[CONTROLLER] = true
	COMPONENTS[MIDDLEWARE] = true
	COMPONENTS[DATABASE] = true
	//COMPONENTS[ROUTER] = true

	ComponentTypes[reflect.TypeOf(BootStrap{})] = true
	ComponentTypes[reflect.TypeOf(Component{})] = true
	ComponentTypes[reflect.TypeOf(Service{})] = true
	ComponentTypes[reflect.TypeOf(Repository{})] = true
	ComponentTypes[reflect.TypeOf(Parameter{})] = true
	ComponentTypes[reflect.TypeOf(Controller{})] = true
	ComponentTypes[reflect.TypeOf(Middleware{})] = true
	ComponentTypes[reflect.TypeOf(Database{})] = true
	//ComponentTypes[reflect.TypeOf(Router{})] = true

	for value, str := range MethodToStr {
		MethodsTypeSet[reflect.TypeOf(value)] = str
	}
}

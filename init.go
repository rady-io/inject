package rhapsody

import "reflect"

func init() {
	COMPONENTS[COMPONENT] = true
	COMPONENTS[SERVICE] = true
	COMPONENTS[REPOSITORY] = true
	//COMPONENTS[CONTROLLER] = true
	//COMPONENTS[ROUTER] = true
	//COMPONENTS[MIDDLEWARE] = true

	ComponentTypes[reflect.TypeOf(Component{})] = true
	ComponentTypes[reflect.TypeOf(Service{})] = true
	ComponentTypes[reflect.TypeOf(Repository{})] = true
	ComponentTypes[reflect.TypeOf(Parameter{})] = true
	//ComponentTypes[reflect.TypeOf(Controller{})] = true
	//ComponentTypes[reflect.TypeOf(Router{})] = true
	//ComponentTypes[reflect.TypeOf(Middleware{})] = true
}

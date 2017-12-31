package types

import "reflect"

func init() {
	COMPONENTS[COMPONENT] = true
	COMPONENTS[SERVICE] = true
	COMPONENTS[CONTROLLER] = true
	COMPONENTS[ROUTER] = true
	COMPONENTS[MIDDLEWARE] = true

	COMPONENT_TYPES[reflect.TypeOf(Component{})] = true
	COMPONENT_TYPES[reflect.TypeOf(Service{})] = true
	COMPONENT_TYPES[reflect.TypeOf(Controller{})] = true
	COMPONENT_TYPES[reflect.TypeOf(Router{})] = true
	COMPONENT_TYPES[reflect.TypeOf(Middleware{})] = true
}

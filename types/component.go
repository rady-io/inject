package types

import "reflect"

type Component struct {
}

type Configuration struct {
}

type Service struct {
}

type Controller struct {
}

type Router struct {
}

type Middleware struct {
}

const COMPONENT = "component"
const CONFIGURATION = "configuration"
const SERVICE = "service"
const CONTROLLER = "controller"
const ROUTER = "router"
const MIDDLEWARE = "middleware"

var COMPONENTS = make(map[string]bool)
var COMPONENT_TYPES = make(map[reflect.Type]bool)



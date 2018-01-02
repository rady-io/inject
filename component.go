package rhapsody

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

type Handler struct {
}

type Repository struct {
} 

const COMPONENT = "component"
const CONFIGURATION = "configuration"
const SERVICE = "service"
const CONTROLLER = "controller"
const ROUTER = "router"
const MIDDLEWARE = "middleware"
const HANDLER = "handler"
const REPOSITORY = "repository"

var COMPONENTS = make(map[string]bool)
var COMPONENT_TYPES = make(map[reflect.Type]bool)



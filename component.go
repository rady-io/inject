package rady

import "reflect"

// BootStrap is a tag to mark a struct as a Bootstrap
type BootStrap struct {
}

// Component is a tag to mark a struct as a Component
type Component struct {
}

// Configuration is a tag to mark a struct as a Configuration
type Configuration struct {
}

// Service is a tag to mark a struct as a Service
type Service struct {
}

// Controller is a tag to mark a struct as a Controller
type Controller struct {
}

// Router is a tag to mark a struct as a Router
type Router struct {
}

// Middleware is a tag to mark a struct as a Middleware
type Middleware struct {
}

// Handler is a tag to mark a struct as a Handler
type Handler struct {
}

// Repository is a tag to mark a struct as a Repository
type Repository struct {
}

// COMPONENT is a tag to mark a field as a Component
const COMPONENT = "component"

// CONFIGURATION is a tag to mark a field as a Configuration
const CONFIGURATION = "configuration"

// SERVICE is a tag to mark a field as a Service
const SERVICE = "service"

// CONTROLLER is a tag to mark a field as a Controller
const CONTROLLER = "controller"

// ROUTER is a tag to mark a field as a Router
const ROUTER = "router"

// MIDDLEWARE is a tag to mark a field as a Middleware
const MIDDLEWARE = "middleware"

// HANDLER is a tag to mark a field as a Handler
const HANDLER = "handler"

// REPOSITORY is a tag to mark a field as a Repository
const REPOSITORY = "repository"

// COMPONENTS is a map to check if a field is COMPONENT, SERVICE or REPOSITORY
var COMPONENTS = make(map[string]bool)

// ComponentTypes is a map to check if a struct is Component, Service, Repository or Parameter
var ComponentTypes = make(map[reflect.Type]bool)



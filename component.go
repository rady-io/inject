package rady

import "reflect"

type (
	// BootStrap is a tag to mark a struct as a Bootstrap
	BootStrap struct {
	}

	// Component is a tag to mark a struct as a Component
	Component struct {
	}

	// Configuration is a tag to mark a struct as a Configuration
	Configuration struct {
	}

	// Service is a tag to mark a struct as a Service
	Service struct {
	}

	// Controller is a tag to mark a struct as a Controller
	Controller struct {
	}

	// Router is a tag to mark a struct as a Router
	Router struct {
	}

	// Middleware is a tag to mark a struct as a Middleware
	Middleware struct {
	}

	// Handler is a tag to mark a struct as a Handler
	Handler struct {
	}

	// Repository is a tag to mark a struct as a Repository
	Repository struct {
	}
)

const (
	// COMPONENT is a tag to mark a field as a Component
	COMPONENT = "component"

	// CONFIGURATION is a tag to mark a field as a Configuration
	CONFIGURATION = "configuration"

	// SERVICE is a tag to mark a field as a Service
	SERVICE = "service"

	// CONTROLLER is a tag to mark a field as a Controller
	CONTROLLER = "controller"

	// ROUTER is a tag to mark a field as a Router
	ROUTER = "router"

	// MIDDLEWARE is a tag to mark a field as a Middleware
	MIDDLEWARE = "middleware"

	// HANDLER is a tag to mark a field as a Handler
	HANDLER = "handler"

	// REPOSITORY is a tag to mark a field as a Repository
	REPOSITORY = "repository"
)

var (
	// COMPONENTS is a map to check if a field is COMPONENT, SERVICE or REPOSITORY
	COMPONENTS = make(map[string]bool)

	// ComponentTypes is a map to check if a struct is Component, Service, Repository or Parameter
	ComponentTypes = make(map[reflect.Type]bool)
)

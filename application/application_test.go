package application

import (
	"summer/types"
	"testing"
	"fmt"
)

type RouterConfig struct {
	types.Configuration
	*UserController
}

type UserController struct {
	types.Component
}

type App struct {
	*RouterConfig
}

func TestCreateApplication(t *testing.T) {
	app := CreateApplication(new(App))
	app.Run()
	for Type, valueMap := range app.BeanMap {
		fmt.Printf("Type: %s\n", Type.String())
		for name, value := range valueMap {
			fmt.Printf("> %s -> %s\n", name, value)
		}
	}
}
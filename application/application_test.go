package application

import (
	"summer/types"
	"testing"
	"fmt"
)

type RouterConfig struct {
	types.Configuration
	*UserComponent `name:"*UserComponent"`
}

type UserComponent struct {
	types.Component
	*RouterConfig
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
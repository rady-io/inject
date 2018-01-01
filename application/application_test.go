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
	types.Controller
}

type App struct {
	*RouterConfig
}

func TestCreateApplication(t *testing.T) {
	app := CreateApplication(new(App))
	app.Run()
	fmt.Println(len(app.BeanMap))
}
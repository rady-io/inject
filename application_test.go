package rady

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

type RouterConfig struct {
	Configuration
	*UserComponent `name:"*UserComponent"`
}

type HandlerConfig struct {
	Configuration
}

type GetUserParam struct {
	Parameter
	*BookService
}

func (rc *RouterConfig) GetUserComponent(BR *BookRepository, GP *GetUserParam) *UserComponent {
	return new(UserComponent)
}

type UserComponent struct {
	Component
	*RouterConfig
}

type BookService struct {
	Service
}

type BookRepository struct {
	Repository
}

type BookController struct {
	Controller `prefix:"/api/v1"`
}

func (b *BookController) GetBooks(ctx *Context) error {
	return nil
}

type App struct {
	*RouterConfig
	*HandlerConfig
	*BookController
}

func TestCreateApplication(t *testing.T) {
	app := CreateApplication(new(App))
	app.Run()
	for Type, valueMap := range app.BeanMap {
		t.Logf("Type: %s\n", Type.String())
		for name, value := range valueMap {
			t.Logf(" %s - Value canset: %t\n", name, value.Value.CanSet())
			t.Logf(" %s - Field canset: %t\n", name, value.Value.Field(0).CanSet())
			assert.True(t, value.Value.Field(0).CanSet(), "Field of %s should CanSet", name)
			assert.True(t, value.Value.CanSet(), "%s should CanSet", name)
		}
	}
}

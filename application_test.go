package rady

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
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
	RedisHost *string `value:"rhapsody.redis.host"`
}

func (rc *RouterConfig) GetUserComponent(BR *BookRepository, GP *GetUserParam) *UserComponent {
	return new(UserComponent)
}

type UserComponent struct {
	Component
	*RouterConfig
	RedisPort *int64 `value:"rhapsody.redis.port"`
}

type BookService struct {
	Service
	RedisPort *int64 `value:"rhapsody.redis.port"`
}

type BookRepository struct {
	Repository
}

type BookController struct {
	Controller `prefix:"/api/v1"`
	GET        `path:"/:id" method:"GetBooks"`
	FILE       `path:"/config" file:"./resources/application.conf"`
	STATIC     `prefix:"/assets" root:"./"`
	BookRepository *BookRepository
}

type BookRouter struct {
	Router `prefix:"/api"`
	*AuthMiddleware
}

type AuthMiddleware struct {
	Middleware
	RedisHost *string `value:"rhapsody.redis.host"`
}

func (a *AuthMiddleware) Auth(next HandlerFunc) HandlerFunc {
	return next
}

func (b *BookController) GetBooks(ctx Context) error {
	return ctx.String(200, fmt.Sprintf(`{"id": "%s"}`, ctx.Param("id")))
}

func (b *BookController) Get0UserUUID(ctx Context) error {
	return ctx.String(200, fmt.Sprintf(`{"uuid": "%s"}`, ctx.Param("uuid")))
}

type App struct {
	*RouterConfig
	*HandlerConfig
	*BookController
	*BookRouter
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

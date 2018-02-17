package rady

import (
	"fmt"
	"testing"
)

type (
	App struct {
		*RouterConfig
		*HandlerConfig
		*BookController
		*BookRouter
	}

	RouterConfig struct {
		Configuration
		*RedisComponent `name:"*RedisComponent"`
	}

	HandlerConfig struct {
		Configuration
	}

	GetUserParam struct {
		Parameter
		*BookService
		RedisHost *string `value:"rady.redis.host" default:"localhost"`
	}

	RedisComponent struct {
		Component
		*RouterConfig
		RedisPort *int64
		RedisHost *string
	}

	BookService struct {
		Service
		//RedisPort *int64 `value:"rady.redis.port"`
	}

	BookRepository struct {
		Repository
		RedisPort *int64 `value:"rady.redis.port" default:"6666"`
	}

	BookController struct {
		Controller `prefix:"/api/v1"`
		GET        `path:"/:id" method:"GetBooks"`
		FILE       `path:"/config" file:"./resources/application.conf"`
		STATIC     `prefix:"/assets" root:"./"`
		BookRepository *BookRepository
		RedisComponent *RedisComponent
		App            *Application
	}

	BookRouter struct {
		Router `prefix:"/api"`
		*AuthMiddleware
	}

	AuthMiddleware struct {
		Middleware
		//RedisHost *string `value:"rady.redis.host"`
	}

	AppTest struct {
		Testing
		*OtherTest
	}

	OtherTest struct {
		Testing
	}
)

func (rc *RouterConfig) GetRedisComponent(BR *BookRepository, GP *GetUserParam) *RedisComponent {
	return &RedisComponent{
		RedisHost: GP.RedisHost,
		RedisPort: BR.RedisPort,
	}
}

func (u *RedisComponent) GetHost() string {
	return *u.RedisHost
}

func (u *RedisComponent) GetPort() int64 {
	return *u.RedisPort
}

func (b *BookController) GetBooks(ctx Context) error {
	return ctx.String(200, fmt.Sprintf(`{"id": "%s"}`, ctx.Param("id")))
}

func (b *BookController) GetUserUUID(ctx Context) error {
	return ctx.String(200, fmt.Sprintf(`{"uuid": "%s"}`, ctx.Param("uuid")))
}

func (b *BookController) GetRedisHost(ctx Context) error {
	return ctx.String(200, fmt.Sprintf(`{"host": "%s"}`, b.RedisComponent.GetHost()))
}

func (b *BookController) GetConfReload(ctx Context) error {
	b.App.ReloadValues()
	return ctx.String(200, fmt.Sprintf(`{"host": "%s"}`, b.RedisComponent.GetHost()))
}

func TestCreateApplication(t *testing.T) {
	CreateTest(new(App)).AddTest(new(AppTest)).AddTests(new(AppTest)).Test(t)
	go CreateApplication(new(App)).Run()
}

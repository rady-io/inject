package rady

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
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
	RedisHost *string `value:"rady.redis.host"`
}

func (rc *RouterConfig) GetUserComponent(BR *BookRepository, GP *GetUserParam) *UserComponent {
	return &UserComponent{
		RedisHost: GP.RedisHost,
		RedisPort: BR.RedisPort,
	}
}

type UserComponent struct {
	Component
	*RouterConfig
	RedisPort *int64
	RedisHost *string
}

func (u *UserComponent) GetHost() string {
	return *u.RedisHost
}

type BookService struct {
	Service
	RedisPort *int64 `value:"rady.redis.port"`
}

type BookRepository struct {
	Repository
	RedisPort *int64 `value:"rady.redis.port"`
}

type BookController struct {
	Controller     `prefix:"/api/v1"`
	GET            `path:"/:id" method:"GetBooks"`
	FILE           `path:"/config" file:"./resources/application.conf"`
	STATIC         `prefix:"/assets" root:"./"`
	BookRepository *BookRepository
	UserComponent  *UserComponent
	App            *Application
}

type BookRouter struct {
	Router `prefix:"/api"`
	*AuthMiddleware
}

type AuthMiddleware struct {
	Middleware
	RedisHost *string `value:"rady.redis.host"`
}

func (a *AuthMiddleware) Auth(next HandlerFunc) HandlerFunc {
	return next
}

func (b *BookController) GetBooks(ctx Context) error {
	return ctx.String(200, fmt.Sprintf(`{"id": "%s"}`, ctx.Param("id")))
}

func (b *BookController) GetUserUUID(ctx Context) error {
	return ctx.String(200, fmt.Sprintf(`{"uuid": "%s"}`, ctx.Param("uuid")))
}

func (b *BookController) GetRedisHost(ctx Context) error {
	return ctx.String(200, fmt.Sprintf(`{"host": "%s"}`, b.UserComponent.GetHost()))
}

func (b *BookController) GetConfReload(ctx Context) error {
	b.App.ReloadValues()
	return ctx.String(200, fmt.Sprintf(`{"host": "%s"}`, b.UserComponent.GetHost()))
}

type App struct {
	*RouterConfig
	*HandlerConfig
	*BookController
	*BookRouter
}

type AppTest struct {
	Testing
	*OtherTest
	RedisPort *int64 `value:"rady.redis.port"`
}

type OtherTest struct {
	Testing
	RedisHost *string `value:"rady.redis.host"`
}

func (a *AppTest) TestRedisPort(t *testing.T) {
	assert.Equal(t, int64(6937), *a.RedisPort)
}

func (o *OtherTest) TestRedisHost(t *testing.T) {
	assert.Equal(t, "127.0.0.1", *o.RedisHost)
}

func TestCreateApplication(t *testing.T) {
	CreateApplication(new(App)).RunTest(t, new(AppTest))
}

package rhapsody

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

type App struct {
	*RouterConfig
	*HandlerConfig
}

func TestCreateApplication(t *testing.T) {
	app := CreateApplication(new(App))
	app.Run()
	for Type, valueMap := range app.BeanMap {
		t.Logf("Type: %s\n", Type.String())
		for name, value := range valueMap {
			t.Logf(" %s - Value canset: %s\n", name, value.Value.CanSet())
			t.Logf(" %s - Field canset: %s\n", name, value.Value.Field(0).CanSet())
			if Type.String() == "*application.Application" && Type.String() == name {
				assert.False(t, value.Value.Field(0).CanSet(), "Field of *application.Application should not CanSet")
			} else {
				assert.True(t, value.Value.Field(0).CanSet(), "Field of %s should CanSet", name)
			}

			assert.True(t, value.Value.CanSet(), "%s should CanSet", name)
		}
	}
}

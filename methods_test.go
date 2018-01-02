package rhapsody

import (
	"testing"
	"reflect"
	"github.com/stretchr/testify/assert"
	"fmt"
)

type UserController struct {
	Controller `prefix:"/api/v1/"`
	GET    `method:"GetIds" path:"/:id"`
	POST   `method:"PostIds" path:"/"`
	PUT
}

func TestGetTags(t *testing.T) {
	handlerValue := reflect.ValueOf(UserController{})
	for i := 0; i < handlerValue.NumField(); i++ {
		switch field := handlerValue.Type().Field(i); field.Type {
		case reflect.TypeOf(Controller{}):
			prefix := field.Tag.Get("prefix")
			t.Logf("Router Prefix: %s\n", prefix)
			assert.Equal(t, prefix, "/api/v1/", fmt.Sprintf("Prefix: %s != /api/v1/", prefix))
		case reflect.TypeOf(GET{}):
			path := field.Tag.Get("path")
			method := field.Tag.Get("method")
			assert.Equal(t, path, "/:id", fmt.Sprintf("Path: %s != /:id", path))
			assert.Equal(t, method, "GetIds", fmt.Sprintf("Method: %s != GetIds", method))
			t.Logf("GET: %s -> %s\n", path, method)
		case reflect.TypeOf(POST{}):
			path := field.Tag.Get("path")
			method := field.Tag.Get("method")
			assert.Equal(t, path, "/", fmt.Sprintf("Path: %s != /", path))
			assert.Equal(t, method, "PostIds", fmt.Sprintf("Method: %s != PostIds", method))
			t.Logf("POST: %s -> %s\n", path, method)
		case reflect.TypeOf(PUT{}):
			path := field.Tag.Get("path")
			method := field.Tag.Get("method")
			assert.Equal(t, path, "", fmt.Sprintf(`Path: %s != ""`, path))
			assert.Equal(t, method, "", fmt.Sprintf(`Method: %s != ""`, method))
			t.Logf("PUT: %s -> %s\n", path, method)
		}
	}
}

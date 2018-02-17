package rady

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type People struct {
	Id   uint32
	Name string
}

type Student struct {
	People
	Stuid string
}

type User struct {
	People
	Uid string
}

type PeopleEntities struct {
	Entities
	*People
	*Student
	*User
}

type PeopleRoot struct {
	*PeopleEntities
}

type EntitiesTest struct {
	*Application
}

func (e *EntitiesTest) TestEntities(t *testing.T) {
	assert.Equal(t, 3, len(e.Entities))
	assert.Equal(t, reflect.TypeOf(new(People)), e.Entities[0])
	assert.Equal(t, reflect.TypeOf(new(Student)), e.Entities[1])
	assert.Equal(t, reflect.TypeOf(new(User)), e.Entities[2])
}

func TestEntities(t *testing.T) {
	CreateTest(new(PeopleRoot)).AddTest(new(EntitiesTest)).Test(t)
}

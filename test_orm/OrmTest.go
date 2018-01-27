package test_orm

import (
	"rady"
	"rady/rorm/sqlite3"
	"testing"
	"github.com/stretchr/testify/assert"
)

type OrmTest struct {
	rady.Testing
	Repo *sqlite3.GormSQLiteRepo
}

func (orm *OrmTest) TestCRUD(t *testing.T) {
	DB := orm.Repo.Begin()
	DB.Create(&User{Name:"xixi", Age: 18})
	NewUser := new(User)
	DB.Where(&User{Name:"xixi"}).Find(&NewUser)

	assert.Equal(t, NewUser.Age, 18)
	DB.Rollback()
}
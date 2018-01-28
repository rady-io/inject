package test_orm

import (
	"rady"
	"rady/rorm/sqlite3"
	"testing"
	"github.com/stretchr/testify/assert"
)

type OrmTest struct {
	rady.Testing
	SQLiteRepo *sqlite3.GormSQLiteRepo
}

func (orm *OrmTest) TestSQLiteRepo(t *testing.T) {
	DB := orm.SQLiteRepo.Begin()
	DB.Create(&User{Name:"xixi", Age: 18})
	NewUser := new(User)
	DB.Where(&User{Name:"xixi"}).Find(&NewUser)

	assert.Equal(t, NewUser.Age, 18)
	DB.Rollback()
}
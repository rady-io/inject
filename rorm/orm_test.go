package rorm

import (
	"github.com/Hexilee/rady"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/Hexilee/rady/rorm/sqlite3"
)

type OrmTest struct {
	rady.Testing
	SQLite *sqlite3.GormSQLite
}

func (orm *OrmTest) TestSQLite(t *testing.T) {
	DB := orm.SQLite.Begin()
	DB.Create(&User{Name: "xixi", Age: 18})
	NewUser := new(User)
	DB.Where(&User{Name: "xixi"}).Find(&NewUser)

	assert.Equal(t, NewUser.Age, 18)
	DB.Rollback()
}


func TestOrm(t *testing.T) {
	rady.CreateApplication(new(OrmRoot)).RunTest(t, new(OrmTest))
}

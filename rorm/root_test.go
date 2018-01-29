package rorm

import "github.com/Hexilee/rady/rorm/sqlite3"

type OrmRoot struct {
	*sqlite3.GormSQLiteConfig
	*OrmStorage
}
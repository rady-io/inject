package test_orm

import "rady/rorm/sqlite3"

type OrmRoot struct {
	*sqlite3.GormSQLiteConfig
	*OrmStorage
}

package rorm

import (
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
    "reflect"
	"os"
	ry "rady"
)

type GormSQLiteParameter struct {
	ry.Parameter
	Path *string `value:"rady.sqlite3.path" default:"./rady.db"`
}

type GormSQLiteRepo struct {
	ry.Repository
	Db *gorm.DB
}

func (g *GormConfig) GetAutoMigrateSQLiteDB(params *GormSQLiteParameter) *GormSQLiteRepo {
	db, err := gorm.Open("sqlite3", *params.Path)
	if err != nil {
		g.App.Logger.Critical("Cannot access to sqlite \nError:\n%s", err.Error())
		os.Exit(1)
	}
	for _, entityType := range g.App.Entities {
		db.AutoMigrate(reflect.New(entityType))
	}
	return &GormSQLiteRepo{Db: db}
}


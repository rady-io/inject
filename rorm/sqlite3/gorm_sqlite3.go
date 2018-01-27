package sqlite3

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"reflect"
	"os"
	ry "rady"
)

type (
	GormSQLiteConfig struct {
		ry.Configuration
		App *ry.Application
	}

	GormSQLiteParameter struct {
		ry.Parameter
		Path *string `value:"rady.sqlite3.path" default:"./rady.db"`
	}

	GormSQLiteRepo struct {
		ry.Repository
		Db *gorm.DB
	}
)

func (g *GormSQLiteConfig) GetAutoMigrateSQLiteDB(params *GormSQLiteParameter) *GormSQLiteRepo {
	db, err := gorm.Open("sqlite3", *params.Path)
	if err != nil {
		g.App.Logger.Critical("Cannot access to sqlite \nError:\n%s", err.Error())
		os.Exit(1)
	}
	for _, entityType := range g.App.Entities {
		if entityType.Kind() == reflect.Ptr && entityType.Elem().Kind() == reflect.Struct {
			g.App.Logger.Debug("AutoMigrate: %s", entityType.String())
			db.AutoMigrate(reflect.New(entityType.Elem()).Interface())
		}
	}
	return &GormSQLiteRepo{Db: db}
}

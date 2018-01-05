package rorm

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"reflect"
	"os"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	ry "rady"
)

type GormSQLServerParameter struct {
	ry.Parameter
	Host      *string `value:"rady.mssql.host" default:"127.0.0.1"`
	Port      *string `value:"rady.mssql.port" default:"1433"`
	Database  *string `value:"rady.mssql.database"`
	Username  *string `value:"rady.mssql.username"`
	Password  *string `value:"rady.mssql.password"`
}

type GormSQLServerRepo struct {
	ry.Repository
	Db *gorm.DB
}

func (g *GormConfig) GetAutoMigrateSQLServerDB(params *GormSQLServerParameter) *GormSQLServerRepo {
	db, err := gorm.Open("mssql", fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", *params.Username, *params.Password, *params.Host, *params.Port, *params.Database))
	if err != nil {
		g.App.Logger.Critical("Cannot connect to mssql \nError:\n%s", err.Error())
		os.Exit(1)
	}
	for _, entityType := range g.App.Entities {
		db.AutoMigrate(reflect.New(entityType))
	}
	return &GormSQLServerRepo{Db: db}
}

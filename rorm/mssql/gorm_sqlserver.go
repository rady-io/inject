package mssql

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	"os"
	"github.com/Hexilee/rady"
	"reflect"
	"github.com/Hexilee/rady/rorm/utils"
)

type (
	GormSQLServerConfig struct {
		rady.Configuration
		App *rady.Application
	}

	GormSQLServerParameter struct {
		rady.Parameter
		Host     *string `value:"rady.mssql.host" default:"127.0.0.1"`
		Port     *string `value:"rady.mssql.port" default:"1433"`
		Database *string `value:"rady.mssql.database"`
		Username *string `value:"rady.mssql.username"`
		Password *string `value:"rady.mssql.password"`
	}

	GormSQLServer struct {
		rady.Database
		*gorm.DB
	}
)

func (g *GormSQLServerConfig) GetAutoMigrateSQLServerDB(params *GormSQLServerParameter) *GormSQLServer {
	db, err := gorm.Open("mssql", fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", *params.Username, *params.Password, *params.Host, *params.Port, *params.Database))
	if err != nil {
		g.App.Logger.Critical("Cannot connect to mssql \nError:\n%s", err.Error())
		os.Exit(1)
	}
	for _, entityType := range g.App.Entities {
		if entityType.Kind() == reflect.Ptr && entityType.Elem().Kind() == reflect.Struct {
			g.App.Logger.Debug("AutoMigrate: %s", entityType.String())
			db.AutoMigrate(reflect.New(entityType.Elem()).Interface())
		}
	}
	return &GormSQLServer{DB: utils.SetGormIfAutoRollback(db)}
}
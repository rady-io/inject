package mysql

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
	"github.com/Hexilee/rady"
	"reflect"
	"github.com/Hexilee/rady/rorm/utils"
)

type (
	GormMySQLConfig struct {
		rady.Configuration
		App *rady.Application
	}

	GormMySQLParameter struct {
		rady.Parameter
		Host      *string `value:"rady.mysql.host" default:"127.0.0.1"`
		Port      *string `value:"rady.mysql.port" default:"3306"`
		Database  *string `value:"rady.mysql.database"`
		Username  *string `value:"rady.mysql.username"`
		Password  *string `value:"rady.mysql.password"`
		Charset   *string `value:"rady.mysql.charset" default:"utf8"`
		ParseTime *string `value:"rady.mysql.parseTime" default:"True"`
		Loc       *string `value:"rady.mysql.loc" default:"Local"`
	}

	GormMySQL struct {
		rady.Database
		*gorm.DB
	}
)

func (g *GormMySQLConfig) GetAutoMigrateMySQLDB(params *GormMySQLParameter) *GormMySQL {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s", *params.Username, *params.Password, *params.Host, *params.Port, *params.Database, *params.Charset, *params.ParseTime, *params.Loc))
	if err != nil {
		g.App.Logger.Critical("Cannot connect to mysql \nError:\n%s", err.Error())
		os.Exit(1)
	}
	for _, entityType := range g.App.Entities {
		if entityType.Kind() == reflect.Ptr && entityType.Elem().Kind() == reflect.Struct {
			g.App.Logger.Debug("AutoMigrate: %s", entityType.String())
			db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(reflect.New(entityType.Elem()).Interface())
		}
	}
	return &GormMySQL{DB: utils.SetGormIfAutoRollback(db)}
}

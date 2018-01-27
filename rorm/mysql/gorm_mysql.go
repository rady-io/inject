package mysql

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"reflect"
	"os"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	ry "rady"
)

type (
	GormMySQLConfig struct {
		ry.Configuration
		App *ry.Application
	}

	GormMySQLParameter struct {
		ry.Parameter
		Host      *string `value:"rady.mysql.host" default:"127.0.0.1"`
		Port      *string `value:"rady.mysql.port" default:"3306"`
		Database  *string `value:"rady.mysql.database"`
		Username  *string `value:"rady.mysql.username"`
		Password  *string `value:"rady.mysql.password"`
		Charset   *string `value:"rady.mysql.charset" default:"utf8"`
		ParseTime *string `value:"rady.mysql.parseTime" default:"True"`
		Loc       *string `value:"rady.mysql.loc" default:"Local"`
	}

	GormMySQLRepo struct {
		ry.Repository
		Db *gorm.DB
	}
)

func (g *GormMySQLConfig) GetAutoMigrateMySQLDB(params *GormMySQLParameter) *GormMySQLRepo {
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
	return &GormMySQLRepo{Db: db}
}

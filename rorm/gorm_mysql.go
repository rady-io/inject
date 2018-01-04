package rorm

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"reflect"
	"os"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	ry "rhapsody"
)

type GormMySQLParameter struct {
	ry.Parameter
	Host      *string `value:"rhapsody.mysql.host" default:"127.0.0.1"`
	Port      *string `value:"rhapsody.mysql.port" default:"3306"`
	Database  *string `value:"rhapsody.mysql.database"`
	Username  *string `value:"rhapsody.mysql.username"`
	Password  *string `value:"rhapsody.mysql.password"`
	Charset   *string `value:"rhapsody.mysql.charset" default:"utf8"`
	ParseTime *string `value:"rhapsody.mysql.parseTime" default:"True"`
	Loc       *string `value:"rhapsody.mysql.loc" default:"Local"`
}

type GormMySQLRepo struct {
	ry.Repository
	Db *gorm.DB
}

func (g *GormConfig) GetAutoMigrateMySQLDB(params *GormMySQLParameter) *GormMySQLRepo {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s", *params.Username, *params.Password, *params.Host, *params.Port, *params.Database, *params.Charset, *params.ParseTime, *params.Loc))
	if err != nil {
		g.App.Logger.Critical("Cannot connect to mysql \nError:\n%s", err.Error())
		os.Exit(1)
	}
	for _, entityType := range g.App.Entities {
		db.AutoMigrate(reflect.New(entityType))
	}
	return &GormMySQLRepo{Db: db}
}

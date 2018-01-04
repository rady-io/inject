package rorm

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"reflect"
	"os"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	ry "rhapsody"
)

type GormPostgresParameter struct {
	ry.Parameter
	Host     *string `value:"rhapsody.postgres.host" default:"127.0.0.1"`
	Port     *string `value:"rhapsody.postgres.port" default:"3306"`
	Database *string `value:"rhapsody.postgres.database"`
	Username *string `value:"rhapsody.postgres.username"`
	Password *string `value:"rhapsody.postgres.password"`
	SSLMode  *string `value:"rhapsody.postgres.sslmode" default:"disable"`
}

type GormPostgresRepo struct {
	ry.Repository
	Db *gorm.DB
}

func (g *GormConfig) GetAutoMigratePostgresDB(params *GormPostgresParameter) *GormPostgresRepo {
	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s", *params.Host, *params.Port, *params.Username, *params.SSLMode, *params.Database, *params.Password))
	if err != nil {
		g.App.Logger.Critical("Cannot connect to postgres \nError:\n%s", err.Error())
		os.Exit(1)
	}
	for _, entityType := range g.App.Entities {
		db.AutoMigrate(reflect.New(entityType))
	}
	return &GormPostgresRepo{Db: db}
}

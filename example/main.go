package main

import (
	"rady"
	"rady/rorm/mysql"
)

func main() {
	rady.CreateApplication(new(Root)).Run()
}

type Root struct {
	rady.CONF `path:"resources/application.yaml"`
	*mysql.GormMySQLConfig
	*StorageEntities
}

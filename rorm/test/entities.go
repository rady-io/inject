package test_orm

import (
	"rady"
	"rady/rorm"
)

type (
	OrmStorage struct {
		rady.Entities
		*User
	}

	User struct {
		rorm.Model
		Name string `gorm:"size:50"`
		Age int
	}
)

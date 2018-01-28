package rorm

import (
	"rady"
)

type (
	OrmStorage struct {
		rady.Entities
		*User
	}

	User struct {
		Model
		Name string `gorm:"size:50"`
		Age  int
	}
)

package rorm

import (
	"github.com/Hexilee/rady"
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

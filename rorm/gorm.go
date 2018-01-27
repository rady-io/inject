package rorm

import (
	ry "rady"
	"github.com/jinzhu/gorm"
)

type (
	GormConfig struct {
		ry.Configuration
		App *ry.Application
	}

	Model = gorm.Model
)

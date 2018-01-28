package rorm

import (
	"github.com/jinzhu/gorm"
	ry "rady"
)

type (
	GormConfig struct {
		ry.Configuration
		App *ry.Application
	}

	Model = gorm.Model
)

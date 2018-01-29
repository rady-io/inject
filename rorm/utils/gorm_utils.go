package utils

import (
	"github.com/jinzhu/gorm"
	"github.com/Hexilee/rady"
)

func SetGormIfAutoRollback(db *gorm.DB) *gorm.DB {
	if rady.IsAutoRollback() {
		db = db.Begin()
	}
	return db
}

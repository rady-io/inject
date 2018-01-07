package main

import (
	"rady"
	"rady/rorm"
)

type StorageEntities struct {
	rady.Entities
	*Org
	*User
	*File
}

type Org struct {
	rorm.Model
	Name  string `gorm:"size:50"`
	Users []User `gorm:"many2many:org_users;"`
	Files []File
}

type User struct {
	rorm.Model
	Name  string `gorm:"size:50"`
	Files []File
}

type File struct {
	rorm.Model
	Hash   string `gorm:"size:32"`
	OrgID  uint
	UserID uint
}

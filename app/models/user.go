package models

import (
	"github.com/goravel/framework/database/orm"
)

type User struct {
	orm.Model
	Name     string `gorm:"column:name"`
	Email    string `gorm:"column:email;uniqueIndex"`
	Password string `gorm:"column:password"`
	Role     string `gorm:"column:role;default:client"` // admin, client
	Admin    *Admin    `gorm:"foreignKey:UserID"`
	Customer *Customer `gorm:"foreignKey:UserID"`
}

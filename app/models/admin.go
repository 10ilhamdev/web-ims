package models

import (
	"github.com/goravel/framework/database/orm"
)

type Admin struct {
	orm.Model
	UserID     uint   `gorm:"column:user_id"`
	Phone      string `gorm:"column:phone"`
	Department string `gorm:"column:department"`
	User       *User  `gorm:"foreignKey:UserID"`
}

package models

import (
	"github.com/goravel/framework/database/orm"
)

type ActivityLog struct {
	orm.Model
	UserID   uint   `gorm:"column:user_id"`
	Activity string `gorm:"column:activity"`
	Device   string `gorm:"column:device"`
	IP       string `gorm:"column:ip"`

	// Relationship
	User *User `gorm:"foreignKey:UserID"`
}

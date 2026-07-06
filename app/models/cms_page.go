package models

import (
	"github.com/goravel/framework/database/orm"
)

type CmsPage struct {
	orm.Model
	Name  string `gorm:"column:name"`
	Type  string `gorm:"column:type"`
	Order int    `gorm:"column:order"`

	// Relationships
	Contents []GuestContent `gorm:"foreignKey:PageID"`
}

  package models

import (
	"github.com/goravel/framework/database/orm"
)

type GuestContent struct {
	orm.Model
	PageID  uint   `gorm:"column:page_id"`
	Key     string `gorm:"column:key;uniqueIndex"`
	ValueId string `gorm:"column:value_id"`
	ValueEn string `gorm:"column:value_en"`
	Section string `gorm:"column:section"`
	Style   string `gorm:"column:style"`
}

package models

import (
	"github.com/goravel/framework/database/orm"
)

type Role struct {
	orm.Model
	Name      string `gorm:"column:name;uniqueIndex"`
	TableName string `gorm:"column:table_name"`
	ModelName string `gorm:"column:model_name"`
	Fields    string `gorm:"column:fields;type:text"`    // JSON array of field definitions
	Relations string `gorm:"column:relations;type:text"` // JSON array of relationship definitions
}

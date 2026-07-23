package models

import (
	"github.com/goravel/framework/database/orm"
)

type Role struct {
	orm.Model
	Name           string `gorm:"column:name;uniqueIndex"`
	Label          string `gorm:"column:label"`
	TableName      string `gorm:"column:table_name"`
	RelationName   string `gorm:"column:relation_name"`
	IsSystem       bool   `gorm:"column:is_system;default:false"`
	IsRegisterable bool   `gorm:"column:is_registerable;default:false"`
	BadgeColor     string `gorm:"column:badge_color"`
	Description    string `gorm:"column:description"`
	DashboardRoute string `gorm:"column:dashboard_route"`
	DashboardView  string `gorm:"column:dashboard_view"`
}

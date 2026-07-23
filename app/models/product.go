package models

import (
	"github.com/goravel/framework/database/orm"
)

type Product struct {
	orm.Model
	Name        string  `gorm:"column:name"`
	Description string  `gorm:"column:description"`
	Price         float64 `gorm:"column:price"`
	OriginalPrice float64 `gorm:"column:original_price"`
	Discount      float64 `gorm:"column:discount"`
	Features      string  `gorm:"column:features"` // JSON string representation of features list
	Image         string  `gorm:"column:image"`    // Icon name or image URL
}

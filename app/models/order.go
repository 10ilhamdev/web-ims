package models

import (
	"github.com/goravel/framework/database/orm"
)

type Order struct {
	orm.Model
	UserID       uint    `gorm:"column:user_id"`
	User         *User   `gorm:"foreignKey:UserID"`
	ProductID    uint    `gorm:"column:product_id"`
	Product      *Product `gorm:"foreignKey:ProductID"`
	Requirements string  `gorm:"column:requirements"`
	Price        float64 `gorm:"column:price"`
	Status       string  `gorm:"column:status;default:pending"` // pending, in_progress, completed, cancelled
}

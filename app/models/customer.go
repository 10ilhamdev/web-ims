package models

import (
	"github.com/goravel/framework/database/orm"
)

type Customer struct {
	orm.Model
	UserID      uint   `gorm:"column:user_id"`
	Phone       string `gorm:"column:phone"`
	CompanyName string `gorm:"column:company_name"`
	Address     string `gorm:"column:address"`
	User        *User  `gorm:"foreignKey:UserID"`
}

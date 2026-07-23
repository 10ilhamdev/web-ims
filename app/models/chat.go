package models

import (
	"github.com/goravel/framework/database/orm"
)

type Chat struct {
	orm.Model
	Token    string        `gorm:"column:token;uniqueIndex"`
	Name     string        `gorm:"column:name"`
	Email    string        `gorm:"column:email"`
	UserID   *uint         `gorm:"column:user_id"`
	Messages []ChatMessage `gorm:"foreignKey:ChatID"`
}

type ChatMessage struct {
	orm.Model
	ChatID     uint   `gorm:"column:chat_id"`
	SenderType string `gorm:"column:sender_type"` // "user" or "admin"
	Message    string `gorm:"column:message"`
}

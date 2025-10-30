package models

import (
	"time"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"size:100;not null;uniqueIndex"`
	Username  string    `json:"username" gorm:"size:100;not null;uniqueIndex"`
	Password  string    `json:"-" gorm:"size:255;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

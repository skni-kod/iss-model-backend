package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type StringArray []string

func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = []string{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, s)
}

func (s StringArray) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "[]", nil
	}
	return json.Marshal(s)
}

type Post struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"size:255;not null"`
	Content     string         `json:"content" gorm:"type:text;not null"`
	Excerpt     string         `json:"excerpt" gorm:"type:text"`
	Author      string         `json:"author" gorm:"size:100"`
	PublishDate string         `json:"publishDate" gorm:"size:100"`
	ReadTime    string         `json:"readTime" gorm:"size:50"`
	Image       string         `json:"image" gorm:"size:500"`
	Images      StringArray    `json:"images" gorm:"type:jsonb"`
	Tags        StringArray    `json:"tags" gorm:"type:jsonb"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Post) TableName() string {
	return "posts"
}

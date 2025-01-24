package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title    string `json:"title"`
	Content  string `json:"content"`
	Category []string `json:"category"` //Can have multiple categories
	UserID   uint   `json:"-"`
	User     User   `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Datetime string `json:"date_time"`
	Score    uint   `json:"score"`
	Picture byte `json:picture`
	Comments []Comment `json:comments`
}

package models

import (
	"gorm.io/gorm"
)


type Comment struct {
    gorm.Model
    Content string `json:"content"`
    UserID  uint   `json:"user_id"` // Matches User model's primary key
    User    User   `json:"user" gorm:"foreignKey:UserID"`
    PostID  uint   `json:"post_id"` // Foreign key for Post
}
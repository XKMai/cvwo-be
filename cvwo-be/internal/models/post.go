package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)


type Post struct {
    gorm.Model
    Title     string         `json:"title"`
    Content   string         `json:"content"`
    Category  pq.StringArray `json:"category" gorm:"type:text[]"`
    UserID    uint           `json:"user_id"`
    User      User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
    Score     uint           `json:"score"`
    Picture   string         `json:"picture"` // URL for image
    Comments  []Comment      `json:"comments" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:PostID"`
}
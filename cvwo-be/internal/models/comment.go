package models

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	Content  string `json:"content"`
	UserId string  `json:"owner_id"`
	User    User   `json:"owner"`
	Post	Post   `json:"post"`
	Datetime string `json:"date_time"`
	Score    uint   `json:"score"`
}

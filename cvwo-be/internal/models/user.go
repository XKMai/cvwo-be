package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Role     string `json:"role"`
	Name     string `json:"name" gorm:"unique"` 
	Password string `json:"password"`
	Description string `json:"description"`
}

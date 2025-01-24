package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Role     string `json:"role"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Picture  byte `json:picture`
	Description string `json:description`
}

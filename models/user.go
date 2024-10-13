package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name  string `json:"name"` // first letter must be Uppercase
	Email string `json:"email"`
}

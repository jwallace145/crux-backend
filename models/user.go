package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Email     string `gorm:"size:100;uniqueIndex;not null" json:"email"`
	FirstName string `gorm:"size:100" json:"first_name"` // optional
	LastName  string `gorm:"size:100" json:"last_name"`  // optional
}

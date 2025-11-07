package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username          string `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Email             string `gorm:"size:100;uniqueIndex;not null" json:"email"`
	PasswordHash      string `gorm:"size:255;not null" json:"-"`          // bcrypt hash, excluded from JSON
	FirstName         string `gorm:"size:100" json:"first_name"`          // optional
	LastName          string `gorm:"size:100" json:"last_name"`           // optional
	ProfilePictureURI string `gorm:"size:255" json:"profile_picture_uri"` // S3 URI for profile picture
}

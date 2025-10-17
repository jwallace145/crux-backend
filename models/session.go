package models

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	User      User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	SessionID string    `gorm:"size:36;uniqueIndex;not null" json:"session_id"`
	IP        string    `gorm:"size:45" json:"ip"`
	UserAgent string    `gorm:"size:256" json:"user_agent"`
	Device    string    `gorm:"size:100" json:"device"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	Revoked   bool      `gorm:"default:false" json:"revoked"`
}

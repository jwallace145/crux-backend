package models

import (
	"gorm.io/gorm"
)

// Route represents a climbing route
type Route struct {
	gorm.Model
	Name        string   `gorm:"size:150;not null" json:"name"`
	Description string   `gorm:"type:text" json:"description,omitempty"`
	Grade       string   `gorm:"size:20;not null" json:"grade"` // Yosemite Decimal System or canonical grade
	Type        string   `gorm:"size:50" json:"type,omitempty"` // e.g., "Sport", "Trad", "Boulder"
	Stars       float32  `json:"stars,omitempty"`               // Optional rating
	NumRatings  int      `json:"num_ratings,omitempty"`         // Optional number of ratings
	Latitude    *float64 `json:"latitude,omitempty"`            // Optional GPS
	Longitude   *float64 `json:"longitude,omitempty"`           // Optional GPS
	WallID      uint     `gorm:"not null;index" json:"wall_id"` // FK to Wall
	Wall        Wall     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
}

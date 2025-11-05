package models

import (
	"gorm.io/gorm"
)

// OutdoorBoulder represents a real-world boulder problem that exists at a climbing area.
// Unlike IndoorBoulder, these are permanent features that are cataloged in the database
// and can be climbed by multiple users over time.
type OutdoorBoulder struct {
	gorm.Model

	// Link to wall/area where boulder is located
	WallID uint `gorm:"not null;index" json:"wall_id"`
	Wall   Wall `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`

	// Boulder details
	Name        string `gorm:"size:150;not null" json:"name"`          // Name of the boulder problem
	Grade       string `gorm:"size:20;not null" json:"grade"`          // V scale (e.g., "V3", "V5")
	Description string `gorm:"type:text" json:"description,omitempty"` // Detailed description, beta, etc.

	// Future enhancements could include:
	// - Multiple grades (sit start, different variations)
	// - FA (First Ascent) information
	// - Star ratings
	// - Reference to climb logs
}

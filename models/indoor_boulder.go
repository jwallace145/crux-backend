package models

import (
	"gorm.io/gorm"
)

// IndoorBoulder outcome constants
const (
	IndoorBoulderOutcomeFell     = "Fell"
	IndoorBoulderOutcomeFlash    = "Flash"
	IndoorBoulderOutcomeOnSight  = "Onsite"
	IndoorBoulderOutcomeRedpoint = "Redpoint"
)

// IndoorBoulder represents a temporary boulder problem at a climbing gym during a training session.
// These are ephemeral problems that get reset/changed by route setters.
// For permanent, real-world boulders, see OutdoorBoulder.
type IndoorBoulder struct {
	gorm.Model

	// Link to training session
	TrainingSessionID uint            `gorm:"not null;index" json:"training_session_id"`
	TrainingSession   TrainingSession `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`

	// Boulder details
	Grade    string  `gorm:"size:20;not null" json:"grade"`      // V scale (e.g., "V3", "V5")
	ColorTag *string `gorm:"size:50" json:"color_tag,omitempty"` // Optional color tag (e.g., "Blue", "Red")
	Outcome  string  `gorm:"size:20;not null" json:"outcome"`    // Fell, Flash, Onsite, Redpoint

	// Optional notes for this specific boulder
	Notes string `gorm:"type:text" json:"notes,omitempty"`
}

// IsSent returns true if the indoor boulder was sent (Flash, Onsite, or Redpoint)
func (ib *IndoorBoulder) IsSent() bool {
	return ib.Outcome == IndoorBoulderOutcomeFlash ||
		ib.Outcome == IndoorBoulderOutcomeOnSight ||
		ib.Outcome == IndoorBoulderOutcomeRedpoint
}

// IsAttempted returns true if the indoor boulder was attempted but not sent (Fell)
func (ib *IndoorBoulder) IsAttempted() bool {
	return ib.Outcome == IndoorBoulderOutcomeFell
}

// HasColorTag returns true if the indoor boulder has a color tag
func (ib *IndoorBoulder) HasColorTag() bool {
	return ib.ColorTag != nil && *ib.ColorTag != ""
}

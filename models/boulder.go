package models

import (
	"gorm.io/gorm"
)

// Boulder outcome constants
const (
	BoulderOutcomeFell     = "Fell"
	BoulderOutcomeFlash    = "Flash"
	BoulderOutcomeOnSight  = "Onsite"
	BoulderOutcomeRedpoint = "Redpoint"
)

// Boulder represents a boulder problem in a training session
type Boulder struct {
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

// IsSent returns true if the boulder was sent (Flash, Onsite, or Redpoint)
func (b *Boulder) IsSent() bool {
	return b.Outcome == BoulderOutcomeFlash ||
		b.Outcome == BoulderOutcomeOnSight ||
		b.Outcome == BoulderOutcomeRedpoint
}

// IsAttempted returns true if the boulder was attempted but not sent (Fell)
func (b *Boulder) IsAttempted() bool {
	return b.Outcome == BoulderOutcomeFell
}

// HasColorTag returns true if the boulder has a color tag
func (b *Boulder) HasColorTag() bool {
	return b.ColorTag != nil && *b.ColorTag != ""
}

package models

import (
	"gorm.io/gorm"
)

// Rope climb type constants
const (
	RopeClimbTypeTopRope = "TR"
	RopeClimbTypeLead    = "Lead"
)

// Rope climb outcome constants
const (
	RopeClimbOutcomeFell     = "Fell"
	RopeClimbOutcomeHung     = "Hung"
	RopeClimbOutcomeFlash    = "Flash"
	RopeClimbOutcomeOnSight  = "Onsite"
	RopeClimbOutcomeRedpoint = "Redpoint"
)

// RopeClimb represents a roped climb (top-rope or lead) in a training session
type RopeClimb struct {
	gorm.Model

	// Link to training session
	TrainingSessionID uint            `gorm:"not null;index" json:"training_session_id"`
	TrainingSession   TrainingSession `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`

	// Climb details
	ClimbType string `gorm:"size:20;not null" json:"climb_type"` // TR or Lead
	Grade     string `gorm:"size:20;not null" json:"grade"`      // YDS scale (e.g., "5.10a", "5.11d")
	Outcome   string `gorm:"size:20;not null" json:"outcome"`    // Fell, Hung, Flash, Onsite, Redpoint

	// Optional notes for this specific climb
	Notes string `gorm:"type:text" json:"notes,omitempty"`
}

// IsTopRope returns true if the climb is top-roped
func (rc *RopeClimb) IsTopRope() bool {
	return rc.ClimbType == RopeClimbTypeTopRope
}

// IsLead returns true if the climb is lead
func (rc *RopeClimb) IsLead() bool {
	return rc.ClimbType == RopeClimbTypeLead
}

// IsSent returns true if the climb was sent (Flash, Onsite, or Redpoint)
func (rc *RopeClimb) IsSent() bool {
	return rc.Outcome == RopeClimbOutcomeFlash ||
		rc.Outcome == RopeClimbOutcomeOnSight ||
		rc.Outcome == RopeClimbOutcomeRedpoint
}

// IsAttempted returns true if the climb was attempted but not sent (Fell or Hung)
func (rc *RopeClimb) IsAttempted() bool {
	return rc.Outcome == RopeClimbOutcomeFell || rc.Outcome == RopeClimbOutcomeHung
}

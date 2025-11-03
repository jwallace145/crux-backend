package models

import (
	"time"

	"gorm.io/gorm"
)

// TrainingSession represents a gym climbing training session
type TrainingSession struct {
	gorm.Model

	// User who logged the session
	UserID uint `gorm:"not null;index" json:"user_id"`
	User   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`

	// Gym where the session took place
	GymID uint `gorm:"not null;index" json:"gym_id"`
	Gym   Gym  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"gym,omitempty"`

	// Session details
	SessionDate time.Time `gorm:"not null;index" json:"session_date"`
	Description string    `gorm:"type:text" json:"description,omitempty"`

	// Training partners (many-to-many relationship)
	Partners []User `gorm:"many2many:training_session_partners;" json:"partners,omitempty"`

	// Climbs during this session
	RopeClimbs []RopeClimb `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"rope_climbs,omitempty"`
	Boulders   []Boulder   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"boulders,omitempty"`
}

// GetTotalClimbs returns the total number of climbs (ropes + boulders) in the session
func (ts *TrainingSession) GetTotalClimbs() int {
	return len(ts.RopeClimbs) + len(ts.Boulders)
}

// GetTotalSends returns the total number of sends (both rope and boulder)
func (ts *TrainingSession) GetTotalSends() int {
	sends := 0
	for _, rc := range ts.RopeClimbs {
		if rc.IsSent() {
			sends++
		}
	}
	for _, b := range ts.Boulders {
		if b.IsSent() {
			sends++
		}
	}
	return sends
}

// HasPartners returns true if the session has training partners
func (ts *TrainingSession) HasPartners() bool {
	return len(ts.Partners) > 0
}

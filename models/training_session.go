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
	RopeClimbs     []RopeClimb     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"rope_climbs,omitempty"`
	IndoorBoulders []IndoorBoulder `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"indoor_boulders,omitempty"`
}

// GetTotalClimbs returns the total number of climbs (ropes + indoor boulders) in the session
func (ts *TrainingSession) GetTotalClimbs() int {
	return len(ts.RopeClimbs) + len(ts.IndoorBoulders)
}

// GetTotalSends returns the total number of sends (both rope and indoor boulder)
func (ts *TrainingSession) GetTotalSends() int {
	sends := 0
	for _, rc := range ts.RopeClimbs {
		if rc.IsSent() {
			sends++
		}
	}
	for _, ib := range ts.IndoorBoulders {
		if ib.IsSent() {
			sends++
		}
	}
	return sends
}

// HasPartners returns true if the session has training partners
func (ts *TrainingSession) HasPartners() bool {
	return len(ts.Partners) > 0
}

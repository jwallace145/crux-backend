package models

import (
	"time"

	"gorm.io/gorm"
)

// ClimbType constants for indoor vs outdoor climbs
const (
	ClimbTypeIndoor  = "indoor"
	ClimbTypeOutdoor = "outdoor"
)

// ClimbStyle constants for different climbing styles
const (
	ClimbStyleLead      = "lead"
	ClimbStyleTopRope   = "top_rope"
	ClimbStyleBoulder   = "boulder"
	ClimbStyleTrad      = "trad"
	ClimbStyleSport     = "sport"
	ClimbStyleAutoBelay = "auto_belay"
)

// Climb represents a user's climbing activity/log entry
// Supports both indoor gym climbs and outdoor route climbs
type Climb struct {
	gorm.Model

	// User relationship - who logged this climb
	UserID uint `gorm:"not null;index" json:"user_id"`
	User   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`

	// Route relationship - optional, only for outdoor climbs linked to a specific route
	RouteID *uint  `gorm:"index" json:"route_id,omitempty"`
	Route   *Route `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"route,omitempty"`

	// Gym relationship - optional, only for indoor climbs linked to a specific gym
	GymID *uint `gorm:"index" json:"gym_id,omitempty"`
	Gym   *Gym  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"gym,omitempty"`

	// Climb type - indoor or outdoor
	ClimbType string `gorm:"size:20;not null;index" json:"climb_type"` // "indoor" or "outdoor"

	// Climb details
	ClimbDate time.Time `gorm:"not null;index" json:"climb_date"` // When the climb occurred
	Grade     string    `gorm:"size:20;not null" json:"grade"`    // Difficulty grade (YDS, V-scale, etc.)
	Style     string    `gorm:"size:50" json:"style,omitempty"`   // lead, top_rope, boulder, etc.

	// Performance tracking
	Completed bool `gorm:"default:true" json:"completed"`    // Whether the climb was completed (sent/topped)
	Attempts  int  `gorm:"default:1" json:"attempts"`        // Number of attempts to complete
	Falls     int  `gorm:"default:0" json:"falls,omitempty"` // Number of falls during the climb

	// Personal rating and notes
	Rating int    `gorm:"default:0" json:"rating,omitempty"` // Personal rating 1-5 stars
	Notes  string `gorm:"type:text" json:"notes,omitempty"`  // Personal notes about the climb
}

// IsIndoor returns true if this is an indoor gym climb
func (c *Climb) IsIndoor() bool {
	return c.ClimbType == ClimbTypeIndoor
}

// IsOutdoor returns true if this is an outdoor climb
func (c *Climb) IsOutdoor() bool {
	return c.ClimbType == ClimbTypeOutdoor
}

// HasRoute returns true if this climb is linked to a specific route
func (c *Climb) HasRoute() bool {
	return c.RouteID != nil
}

// HasGym returns true if this climb is linked to a specific gym
func (c *Climb) HasGym() bool {
	return c.GymID != nil
}

// GetDisplayName returns an appropriate display name for the climb
func (c *Climb) GetDisplayName() string {
	// If linked to a route, use the route name
	if c.HasRoute() && c.Route != nil {
		return c.Route.Name
	}

	// Default to just the grade
	return c.Grade
}

// GetLocation returns an appropriate location string for the climb
func (c *Climb) GetLocation() string {
	// If linked to a route with a crag, use that
	if c.HasRoute() && c.Route != nil && c.Route.Wall.Crag.Name != "" {
		return c.Route.Wall.Crag.Name
	}
	return "UNKNOWN_LOCATION"
}

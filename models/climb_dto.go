package models

import (
	"time"
)

// CreateClimbRequest represents the request body for creating a new climb
type CreateClimbRequest struct {
	// Required fields
	ClimbType string    `json:"climb_type" validate:"required,oneof=indoor outdoor"`
	ClimbDate time.Time `json:"climb_date" validate:"required"`
	Grade     string    `json:"grade" validate:"required,min=1,max=20"`

	// Optional relationship IDs
	RouteID *uint `json:"route_id,omitempty"` // For outdoor climbs linked to a route
	GymID   *uint `json:"gym_id,omitempty"`   // For indoor climbs linked to a gym

	// Optional details
	Style string `json:"style,omitempty" validate:"omitempty,max=50"`

	// Performance tracking (with defaults)
	Completed bool `json:"completed"`          // Defaults to true
	Attempts  int  `json:"attempts,omitempty"` // Defaults to 1
	Falls     int  `json:"falls,omitempty"`    // Defaults to 0

	// Personal rating and notes
	Rating int    `json:"rating,omitempty" validate:"omitempty,min=0,max=5"`
	Notes  string `json:"notes,omitempty" validate:"omitempty,max=1000"`
}

// ClimbResponse represents the climb data returned in API responses
type ClimbResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	ClimbType string    `json:"climb_type"`
	ClimbDate time.Time `json:"climb_date"`
	Grade     string    `json:"grade"`
	Style     string    `json:"style,omitempty"`

	// Optional relationships
	RouteID *uint `json:"route_id,omitempty"`
	GymID   *uint `json:"gym_id,omitempty"`

	// Performance
	Completed bool `json:"completed"`
	Attempts  int  `json:"attempts"`
	Falls     int  `json:"falls,omitempty"`

	// Personal rating and notes
	Rating int    `json:"rating,omitempty"`
	Notes  string `json:"notes,omitempty"`

	// Metadata
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToClimbResponse converts a Climb model to a ClimbResponse DTO
func (c *Climb) ToClimbResponse() *ClimbResponse {
	return &ClimbResponse{
		ID:        c.ID,
		UserID:    c.UserID,
		ClimbType: c.ClimbType,
		ClimbDate: c.ClimbDate,
		Grade:     c.Grade,
		Style:     c.Style,
		RouteID:   c.RouteID,
		GymID:     c.GymID,
		Completed: c.Completed,
		Attempts:  c.Attempts,
		Falls:     c.Falls,
		Rating:    c.Rating,
		Notes:     c.Notes,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

package models

import (
	"time"
)

// BoulderRequest represents a boulder problem in the create request
type BoulderRequest struct {
	Grade    string  `json:"grade" validate:"required,min=1,max=20"`
	ColorTag *string `json:"color_tag,omitempty" validate:"omitempty,max=50"`
	Outcome  string  `json:"outcome" validate:"required,oneof=Fell Flash Onsite Redpoint"`
	Notes    string  `json:"notes,omitempty" validate:"omitempty,max=1000"`
}

// RopeClimbRequest represents a rope climb in the create request
type RopeClimbRequest struct {
	ClimbType string `json:"climb_type" validate:"required,oneof=TR Lead"`
	Grade     string `json:"grade" validate:"required,min=1,max=20"`
	Outcome   string `json:"outcome" validate:"required,oneof=Fell Hung Flash Onsite Redpoint"`
	Notes     string `json:"notes,omitempty" validate:"omitempty,max=1000"`
}

// CreateTrainingSessionRequest represents the request body for creating a new training session
type CreateTrainingSessionRequest struct {
	// Required fields
	GymID       uint      `json:"gym_id" validate:"required"`
	SessionDate time.Time `json:"session_date" validate:"required"`

	// Optional fields
	Description string `json:"description,omitempty" validate:"omitempty,max=1000"`
	PartnerIDs  []uint `json:"partner_ids,omitempty"`

	// Climbs during the session
	Boulders   []BoulderRequest   `json:"boulders,omitempty"`
	RopeClimbs []RopeClimbRequest `json:"rope_climbs,omitempty"`
}

// BoulderResponse represents a boulder problem in the response
type BoulderResponse struct {
	ID                uint      `json:"id"`
	TrainingSessionID uint      `json:"training_session_id"`
	Grade             string    `json:"grade"`
	ColorTag          *string   `json:"color_tag,omitempty"`
	Outcome           string    `json:"outcome"`
	Notes             string    `json:"notes,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// RopeClimbResponse represents a rope climb in the response
type RopeClimbResponse struct {
	ID                uint      `json:"id"`
	TrainingSessionID uint      `json:"training_session_id"`
	ClimbType         string    `json:"climb_type"`
	Grade             string    `json:"grade"`
	Outcome           string    `json:"outcome"`
	Notes             string    `json:"notes,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// PartnerResponse represents a training partner in the response
type PartnerResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

// GymResponse represents basic gym information in the response
type GymResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	City string `json:"city,omitempty"`
}

// TrainingSessionResponse represents the training session data returned in API responses
type TrainingSessionResponse struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	GymID       uint      `json:"gym_id"`
	SessionDate time.Time `json:"session_date"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Nested relationships
	Gym        *GymResponse        `json:"gym,omitempty"`
	Partners   []PartnerResponse   `json:"partners,omitempty"`
	Boulders   []BoulderResponse   `json:"boulders,omitempty"`
	RopeClimbs []RopeClimbResponse `json:"rope_climbs,omitempty"`

	// Statistics
	TotalClimbs int `json:"total_climbs"`
	TotalSends  int `json:"total_sends"`
}

// ToTrainingSessionResponse converts a TrainingSession model to a TrainingSessionResponse DTO
func (ts *TrainingSession) ToTrainingSessionResponse() *TrainingSessionResponse {
	response := &TrainingSessionResponse{
		ID:          ts.ID,
		UserID:      ts.UserID,
		GymID:       ts.GymID,
		SessionDate: ts.SessionDate,
		Description: ts.Description,
		CreatedAt:   ts.CreatedAt,
		UpdatedAt:   ts.UpdatedAt,
		TotalClimbs: ts.GetTotalClimbs(),
		TotalSends:  ts.GetTotalSends(),
	}

	// Include gym information if loaded
	if ts.Gym.ID != 0 {
		response.Gym = &GymResponse{
			ID:   ts.Gym.ID,
			Name: ts.Gym.Name,
			City: ts.Gym.City,
		}
	}

	// Include partners if loaded
	if len(ts.Partners) > 0 {
		response.Partners = make([]PartnerResponse, len(ts.Partners))
		for i, partner := range ts.Partners {
			response.Partners[i] = PartnerResponse{
				ID:        partner.ID,
				Username:  partner.Username,
				FirstName: partner.FirstName,
				LastName:  partner.LastName,
			}
		}
	}

	// Include boulders if loaded
	if len(ts.Boulders) > 0 {
		response.Boulders = make([]BoulderResponse, len(ts.Boulders))
		for i, boulder := range ts.Boulders {
			response.Boulders[i] = BoulderResponse{
				ID:                boulder.ID,
				TrainingSessionID: boulder.TrainingSessionID,
				Grade:             boulder.Grade,
				ColorTag:          boulder.ColorTag,
				Outcome:           boulder.Outcome,
				Notes:             boulder.Notes,
				CreatedAt:         boulder.CreatedAt,
				UpdatedAt:         boulder.UpdatedAt,
			}
		}
	}

	// Include rope climbs if loaded
	if len(ts.RopeClimbs) > 0 {
		response.RopeClimbs = make([]RopeClimbResponse, len(ts.RopeClimbs))
		for i, ropeClimb := range ts.RopeClimbs {
			response.RopeClimbs[i] = RopeClimbResponse{
				ID:                ropeClimb.ID,
				TrainingSessionID: ropeClimb.TrainingSessionID,
				ClimbType:         ropeClimb.ClimbType,
				Grade:             ropeClimb.Grade,
				Outcome:           ropeClimb.Outcome,
				Notes:             ropeClimb.Notes,
				CreatedAt:         ropeClimb.CreatedAt,
				UpdatedAt:         ropeClimb.UpdatedAt,
			}
		}
	}

	return response
}

// ToBoulder converts a BoulderRequest to a Boulder model
func (br *BoulderRequest) ToBoulder() *Boulder {
	return &Boulder{
		Grade:    br.Grade,
		ColorTag: br.ColorTag,
		Outcome:  br.Outcome,
		Notes:    br.Notes,
	}
}

// ToRopeClimb converts a RopeClimbRequest to a RopeClimb model
func (rcr *RopeClimbRequest) ToRopeClimb() *RopeClimb {
	return &RopeClimb{
		ClimbType: rcr.ClimbType,
		Grade:     rcr.Grade,
		Outcome:   rcr.Outcome,
		Notes:     rcr.Notes,
	}
}

package models

import (
	"time"
)

// CreateGymRequest represents the request body for creating a new gym
type CreateGymRequest struct {
	// Required fields
	Name    string `json:"name" validate:"required,min=1,max=200"`
	Type    string `json:"type" validate:"required,oneof=bouldering roped full"`
	City    string `json:"city" validate:"required,min=1,max=100"`
	Country string `json:"country" validate:"required,min=1,max=100"`

	// Optional basic information
	Description string `json:"description,omitempty" validate:"omitempty,max=5000"`

	// Optional location details
	Address    string `json:"address,omitempty" validate:"omitempty,max=300"`
	State      string `json:"state,omitempty" validate:"omitempty,max=100"`
	Province   string `json:"province,omitempty" validate:"omitempty,max=100"`
	PostalCode string `json:"postal_code,omitempty" validate:"omitempty,max=20"`

	// GPS coordinates
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`

	// Contact information
	Phone   string `json:"phone,omitempty" validate:"omitempty,max=50"`
	Email   string `json:"email,omitempty" validate:"omitempty,email,max=200"`
	Website string `json:"website,omitempty" validate:"omitempty,url,max=300"`

	// Operating hours
	Hours string `json:"hours,omitempty" validate:"omitempty,max=1000"`

	// Facilities and features
	HasBouldering   bool `json:"has_bouldering"`
	HasTopRope      bool `json:"has_top_rope"`
	HasLeadClimbing bool `json:"has_lead_climbing"`
	HasAutoBelay    bool `json:"has_auto_belay"`
	HasKidsArea     bool `json:"has_kids_area"`
	HasTrainingArea bool `json:"has_training_area"`
	HasYogaClasses  bool `json:"has_yoga_classes"`
	HasShower       bool `json:"has_shower"`
	HasParking      bool `json:"has_parking"`
	HasGearRental   bool `json:"has_gear_rental"`
	HasProShop      bool `json:"has_pro_shop"`
	HasCafe         bool `json:"has_cafe"`

	// Capacity and size
	WallHeight int `json:"wall_height,omitempty" validate:"omitempty,min=0"`
	SquareFeet int `json:"square_feet,omitempty" validate:"omitempty,min=0"`

	// Pricing
	DayPassPrice    float64 `json:"day_pass_price,omitempty" validate:"omitempty,min=0"`
	MonthlyPrice    float64 `json:"monthly_price,omitempty" validate:"omitempty,min=0"`
	YearlyPrice     float64 `json:"yearly_price,omitempty" validate:"omitempty,min=0"`
	GearRentalPrice float64 `json:"gear_rental_price,omitempty" validate:"omitempty,min=0"`

	// Additional information
	Notes  string `json:"notes,omitempty" validate:"omitempty,max=5000"`
	Active bool   `json:"active"`
}

// FullGymResponse represents the complete gym data returned in API responses
type FullGymResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type"`

	// Location details
	Address    string `json:"address,omitempty"`
	City       string `json:"city"`
	State      string `json:"state,omitempty"`
	Province   string `json:"province,omitempty"`
	Country    string `json:"country"`
	PostalCode string `json:"postal_code,omitempty"`

	// GPS coordinates
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`

	// Contact information
	Phone   string `json:"phone,omitempty"`
	Email   string `json:"email,omitempty"`
	Website string `json:"website,omitempty"`

	// Operating hours
	Hours string `json:"hours,omitempty"`

	// Facilities and features
	HasBouldering   bool `json:"has_bouldering"`
	HasTopRope      bool `json:"has_top_rope"`
	HasLeadClimbing bool `json:"has_lead_climbing"`
	HasAutoBelay    bool `json:"has_auto_belay"`
	HasKidsArea     bool `json:"has_kids_area"`
	HasTrainingArea bool `json:"has_training_area"`
	HasYogaClasses  bool `json:"has_yoga_classes"`
	HasShower       bool `json:"has_shower"`
	HasParking      bool `json:"has_parking"`
	HasGearRental   bool `json:"has_gear_rental"`
	HasProShop      bool `json:"has_pro_shop"`
	HasCafe         bool `json:"has_cafe"`

	// Capacity and size
	WallHeight int `json:"wall_height,omitempty"`
	SquareFeet int `json:"square_feet,omitempty"`

	// Pricing
	DayPassPrice    float64 `json:"day_pass_price,omitempty"`
	MonthlyPrice    float64 `json:"monthly_price,omitempty"`
	YearlyPrice     float64 `json:"yearly_price,omitempty"`
	GearRentalPrice float64 `json:"gear_rental_price,omitempty"`

	// Additional information
	Notes  string `json:"notes,omitempty"`
	Active bool   `json:"active"`

	// Metadata
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToFullGymResponse converts a Gym model to a FullGymResponse DTO
func (g *Gym) ToFullGymResponse() *FullGymResponse {
	return &FullGymResponse{
		ID:              g.ID,
		Name:            g.Name,
		Description:     g.Description,
		Type:            g.Type,
		Address:         g.Address,
		City:            g.City,
		State:           g.State,
		Province:        g.Province,
		Country:         g.Country,
		PostalCode:      g.PostalCode,
		Latitude:        g.Latitude,
		Longitude:       g.Longitude,
		Phone:           g.Phone,
		Email:           g.Email,
		Website:         g.Website,
		Hours:           g.Hours,
		HasBouldering:   g.HasBouldering,
		HasTopRope:      g.HasTopRope,
		HasLeadClimbing: g.HasLeadClimbing,
		HasAutoBelay:    g.HasAutoBelay,
		HasKidsArea:     g.HasKidsArea,
		HasTrainingArea: g.HasTrainingArea,
		HasYogaClasses:  g.HasYogaClasses,
		HasShower:       g.HasShower,
		HasParking:      g.HasParking,
		HasGearRental:   g.HasGearRental,
		HasProShop:      g.HasProShop,
		HasCafe:         g.HasCafe,
		WallHeight:      g.WallHeight,
		SquareFeet:      g.SquareFeet,
		DayPassPrice:    g.DayPassPrice,
		MonthlyPrice:    g.MonthlyPrice,
		YearlyPrice:     g.YearlyPrice,
		GearRentalPrice: g.GearRentalPrice,
		Notes:           g.Notes,
		Active:          g.Active,
		CreatedAt:       g.CreatedAt,
		UpdatedAt:       g.UpdatedAt,
	}
}

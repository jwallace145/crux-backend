package models

import (
	"gorm.io/gorm"
)

// GymType constants for different types of climbing gyms
const (
	GymTypeBouldering = "bouldering"
	GymTypeRoped      = "roped"
	GymTypeFull       = "full" // Both bouldering and roped climbing
)

// Gym represents a climbing gym/facility
type Gym struct {
	gorm.Model

	// Basic information
	Name        string `gorm:"size:200;not null;index" json:"name"`
	Description string `gorm:"type:text" json:"description,omitempty"`
	Type        string `gorm:"size:50;not null" json:"type"` // bouldering, roped, full

	// Location details
	Address    string `gorm:"size:300" json:"address,omitempty"`
	City       string `gorm:"size:100;not null;index" json:"city"`
	State      string `gorm:"size:100" json:"state,omitempty"`    // For US/Canada
	Province   string `gorm:"size:100" json:"province,omitempty"` // Alternative to State
	Country    string `gorm:"size:100;not null;index" json:"country"`
	PostalCode string `gorm:"size:20" json:"postal_code,omitempty"`

	// GPS coordinates for mapping
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`

	// Contact information
	Phone   string `gorm:"size:50" json:"phone,omitempty"`
	Email   string `gorm:"size:200" json:"email,omitempty"`
	Website string `gorm:"size:300" json:"website,omitempty"`

	// Operating hours (could be enhanced with structured schedule later)
	Hours string `gorm:"type:text" json:"hours,omitempty"` // Free-form text for now

	// Facilities and features
	HasBouldering   bool `gorm:"default:false" json:"has_bouldering"`
	HasTopRope      bool `gorm:"default:false" json:"has_top_rope"`
	HasLeadClimbing bool `gorm:"default:false" json:"has_lead_climbing"`
	HasAutoBelay    bool `gorm:"default:false" json:"has_auto_belay"`
	HasKidsArea     bool `gorm:"default:false" json:"has_kids_area"`
	HasTrainingArea bool `gorm:"default:false" json:"has_training_area"`
	HasYogaClasses  bool `gorm:"default:false" json:"has_yoga_classes"`
	HasShower       bool `gorm:"default:false" json:"has_shower"`
	HasParking      bool `gorm:"default:false" json:"has_parking"`
	HasGearRental   bool `gorm:"default:false" json:"has_gear_rental"`
	HasProShop      bool `gorm:"default:false" json:"has_pro_shop"`
	HasCafe         bool `gorm:"default:false" json:"has_cafe"`

	// Capacity and size
	WallHeight int `json:"wall_height,omitempty"` // Max wall height in feet
	SquareFeet int `json:"square_feet,omitempty"` // Total climbing area

	// Pricing (simplified - could be enhanced with membership tiers later)
	DayPassPrice    float64 `json:"day_pass_price,omitempty"`
	MonthlyPrice    float64 `json:"monthly_price,omitempty"`
	YearlyPrice     float64 `json:"yearly_price,omitempty"`
	GearRentalPrice float64 `json:"gear_rental_price,omitempty"`

	// Additional information
	Notes string `gorm:"type:text" json:"notes,omitempty"` // Additional notes or comments

	// Status
	Active bool `gorm:"default:true" json:"active"` // Whether the gym is currently open/operational
}

// GetFullAddress returns a formatted full address string
func (g *Gym) GetFullAddress() string {
	address := g.Address
	if g.City != "" {
		if address != "" {
			address += ", "
		}
		address += g.City
	}
	if g.State != "" {
		if address != "" {
			address += ", "
		}
		address += g.State
	} else if g.Province != "" {
		if address != "" {
			address += ", "
		}
		address += g.Province
	}
	if g.PostalCode != "" {
		if address != "" {
			address += " "
		}
		address += g.PostalCode
	}
	if g.Country != "" {
		if address != "" {
			address += ", "
		}
		address += g.Country
	}
	return address
}

// GetLocation returns a simplified location string (City, State/Province)
func (g *Gym) GetLocation() string {
	location := g.City
	if g.State != "" {
		if location != "" {
			location += ", "
		}
		location += g.State
	} else if g.Province != "" {
		if location != "" {
			location += ", "
		}
		location += g.Province
	}
	return location
}

// HasRopedClimbing returns true if the gym has any roped climbing facilities
func (g *Gym) HasRopedClimbing() bool {
	return g.HasTopRope || g.HasLeadClimbing || g.HasAutoBelay
}

// GetFacilities returns a list of available facilities
func (g *Gym) GetFacilities() []string {
	facilities := []string{}

	if g.HasBouldering {
		facilities = append(facilities, "Bouldering")
	}
	if g.HasTopRope {
		facilities = append(facilities, "Top Rope")
	}
	if g.HasLeadClimbing {
		facilities = append(facilities, "Lead Climbing")
	}
	if g.HasAutoBelay {
		facilities = append(facilities, "Auto Belay")
	}
	if g.HasKidsArea {
		facilities = append(facilities, "Kids Area")
	}
	if g.HasTrainingArea {
		facilities = append(facilities, "Training Area")
	}
	if g.HasYogaClasses {
		facilities = append(facilities, "Yoga Classes")
	}
	if g.HasShower {
		facilities = append(facilities, "Showers")
	}
	if g.HasParking {
		facilities = append(facilities, "Parking")
	}
	if g.HasGearRental {
		facilities = append(facilities, "Gear Rental")
	}
	if g.HasProShop {
		facilities = append(facilities, "Pro Shop")
	}
	if g.HasCafe {
		facilities = append(facilities, "Cafe")
	}

	return facilities
}

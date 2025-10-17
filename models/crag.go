package models

import (
	"gorm.io/gorm"
)

type Crag struct {
	gorm.Model
	Name      string   `gorm:"size:150;not null" json:"name"`
	State     string   `gorm:"size:50" json:"state,omitempty"`
	City      string   `gorm:"size:100" json:"city,omitempty"`
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
	Walls     []Wall   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"walls,omitempty"`
}

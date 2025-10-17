package models

import (
	"gorm.io/gorm"
)

type Wall struct {
	gorm.Model
	Name   string  `gorm:"size:150;not null" json:"name"`
	CragID uint    `gorm:"not null;index" json:"crag_id"`
	Crag   Crag    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	Routes []Route `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"routes,omitempty"`
}

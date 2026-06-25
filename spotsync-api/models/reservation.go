package models

import "time"

type Reservation struct {
	ID           uint        `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID       uint        `json:"user_id" gorm:"not null;index"`
	ZoneID       uint        `json:"zone_id" gorm:"not null;index"`
	LicensePlate string      `json:"license_plate" gorm:"not null;size:15"`
	Status       string      `json:"status" gorm:"not null;default:'active';index"`
	User         User        `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Zone         ParkingZone `json:"zone,omitempty" gorm:"foreignKey:ZoneID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

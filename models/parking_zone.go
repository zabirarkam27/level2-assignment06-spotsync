package models

import "time"

type ParkingZone struct {
	ID            uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name          string    `json:"name" gorm:"not null"`
	Type          string    `json:"type" gorm:"not null"`
	TotalCapacity int       `json:"total_capacity" gorm:"not null"`
	PricePerHour  float64   `json:"price_per_hour" gorm:"not null"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

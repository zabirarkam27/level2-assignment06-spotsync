package dto

import "time"

type CreateReservationRequest struct {
	ZoneID       uint   `json:"zone_id" validate:"required,gt=0"`
	LicensePlate string `json:"license_plate" validate:"required,max=15"`
}

type ZoneInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type ReservationUserInfo struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type ReservationResponse struct {
	ID           uint                 `json:"id"`
	UserID       *uint                `json:"user_id,omitempty"`
	ZoneID       *uint                `json:"zone_id,omitempty"`
	LicensePlate string               `json:"license_plate"`
	Status       string               `json:"status"`
	User         *ReservationUserInfo `json:"user,omitempty"`
	Zone         *ZoneInfo            `json:"zone,omitempty"`
	CreatedAt    time.Time            `json:"created_at"`
	UpdatedAt    *time.Time           `json:"updated_at,omitempty"`
}

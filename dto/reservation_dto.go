package dto

import "time"

// ---- Reservation DTOs ----

type CreateReservationRequest struct {
	ZoneID       uint   `json:"zone_id" validate:"required"`
	LicensePlate string `json:"license_plate" validate:"required,max=15"`
}

type ZoneInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type ReservationResponse struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id,omitempty"`
	ZoneID       uint      `json:"zone_id,omitempty"`
	LicensePlate string    `json:"license_plate"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}

type MyReservationResponse struct {
	ID           uint      `json:"id"`
	LicensePlate string    `json:"license_plate"`
	Status       string    `json:"status"`
	Zone         ZoneInfo  `json:"zone"`
	CreatedAt    time.Time `json:"created_at"`
}

type AdminReservationResponse struct {
	ID           uint             `json:"id"`
	LicensePlate string           `json:"license_plate"`
	Status       string           `json:"status"`
	User         LoginUserResponse `json:"user"`
	Zone         ZoneInfo         `json:"zone"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
}

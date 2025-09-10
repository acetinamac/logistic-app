package domain

import "time"

// Coordinates table
type Coordinates struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Latitude  float64   `json:"latitude" gorm:"type:decimal(10,8);not null"`
	Longitude float64   `json:"longitude" gorm:"type:decimal(11,8);not null"`
	CreatedAt time.Time `json:"created_at"`
}

// Address table
type Address struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	CustomerID     uint      `json:"customer_id" gorm:"not null"`
	Street         string    `json:"street" gorm:"size:255;not null"`
	ExteriorNumber string    `json:"exterior_number" gorm:"size:10"`
	InteriorNumber string    `json:"interior_number" gorm:"size:10"`
	Neighborhood   string    `json:"neighborhood" gorm:"size:100"`
	PostalCode     string    `json:"postal_code" gorm:"size:10"`
	City           string    `json:"city" gorm:"size:100;not null"`
	State          string    `json:"state" gorm:"size:100;not null"`
	Country        string    `json:"country" gorm:"size:100;default:Mexico;not null"`
	CoordinateID   *uint     `json:"coordinate_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	IsActive       bool      `json:"is_active" gorm:"default:true;not null"`
}

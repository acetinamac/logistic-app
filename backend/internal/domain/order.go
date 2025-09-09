package domain

import "time"

type OrderStatus string

const (
	OrderCreated   OrderStatus = "created"
	OrderCollected OrderStatus = "collected"
	OrderInStation OrderStatus = "in_station"
	OrderInRoute   OrderStatus = "in_route"
	OrderDelivered OrderStatus = "delivered"
	OrderCancelled OrderStatus = "cancelled"
)

type PackageSize string

const (
	PackageS  PackageSize = "S"
	PackageM  PackageSize = "M"
	PackageL  PackageSize = "L"
	PackageXL PackageSize = "XL"
)

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
}

// Package types table
type PackageType struct {
	ID          uint        `json:"id" gorm:"primaryKey"`
	SizeCode    PackageSize `json:"size_code" gorm:"type:package_size_enum;unique;not null"`
	MaxWeightKg float64     `json:"max_weight_kg" gorm:"type:decimal(5,2);not null"`
	Description string      `json:"description" gorm:"type:text"`
	IsActive    bool        `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time   `json:"created_at"`
}

// Orders table
type Order struct {
	ID                   uint        `json:"id" gorm:"primaryKey"`
	OrderNumber          string      `json:"order_number" gorm:"size:50;uniqueIndex"`
	OriginAddressID      uint        `json:"origin_address_id" gorm:"not null"`
	DestinationAddressID uint        `json:"destination_address_id" gorm:"not null"`
	PackageTypeID        uint        `json:"package_type_id" gorm:"not null"`
	ActualWeightKg       float64     `json:"actual_weight_kg" gorm:"type:decimal(5,2)"`
	Status               OrderStatus `json:"status" gorm:"type:order_status_enum;default:created;not null"`
	CustomerID           uint        `json:"customer_id" gorm:"not null"`
	CreatedBy            uint        `json:"created_by" gorm:"not null"`
	UpdatedBy            *uint       `json:"updated_by"`
	CreatedAt            time.Time   `json:"created_at"`
	UpdatedAt            time.Time   `json:"updated_at"`
	Observations         string      `json:"observations" gorm:"type:text"`
	InternalNotes        string      `json:"internal_notes" gorm:"type:text"`
}

// Order status history table
type OrderStatusHistory struct {
	ID             uint        `json:"id" gorm:"primaryKey"`
	OrderID        uint        `json:"order_id" gorm:"not null;index"`
	PreviousStatus OrderStatus `json:"previous_status" gorm:"type:order_status_enum"`
	NewStatus      OrderStatus `json:"new_status" gorm:"type:order_status_enum;not null"`
	ChangedAt      time.Time   `json:"changed_at"`
	ChangedBy      uint        `json:"changed_by" gorm:"not null"`
	Notes          string      `json:"notes" gorm:"type:text"`
}

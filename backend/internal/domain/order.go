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

// Orders table
type Order struct {
	ID                   uint        `json:"id" gorm:"primaryKey"`
	OrderNumber          string      `json:"order_number" gorm:"size:50;uniqueIndex"`
	OriginAddressID      uint        `json:"origin_address_id" gorm:"not null"`
	DestinationAddressID uint        `json:"destination_address_id" gorm:"not null"`
	PackageTypeID        uint        `json:"package_type_id" gorm:"not null"`
	Quantity             uint        `json:"quantity" gorm:"not null"`
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

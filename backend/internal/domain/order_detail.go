package domain

import "time"

// OrderDetail represents the detailed view of an order with joined info
// Mirrors the columns requested in the SQL provided by the user
// JSON tags match snake_case names expected by the frontend

type OrderDetail struct {
	ID                   uint        `json:"id"`
	OrderNumber          string      `json:"order_number"`
	CreatedAt            time.Time   `json:"created_at"`
	UserID               uint        `json:"user_id"`
	FullName             string      `json:"full_name"`
	OriginAddressID      uint        `json:"origin_address_id"`
	AOStreet             string      `json:"ao_street"`
	AOExterior           string      `json:"ao_exterior"`
	AONeighborhood       string      `json:"ao_neighborhood"`
	AOCity               string      `json:"ao_city"`
	AOPostal             string      `json:"ao_postal"`
	DestinationAddressID uint        `json:"destination_address_id"`
	ADStreet             string      `json:"ad_street"`
	ADExterior           string      `json:"ad_exterior"`
	ADNeighborhood       string      `json:"ad_neighborhood"`
	ADCity               string      `json:"ad_city"`
	ADPostal             string      `json:"ad_postal"`
	ActualWeightKg       float64     `json:"actual_weight_kg"`
	PackageTypeID        uint        `json:"package_type_id"`
	SizeCode             PackageSize `json:"size_code"`
	Observations         string      `json:"observations"`
	InternalNotes        string      `json:"internal_notes"`
	UpdatedAt            time.Time   `json:"updated_at"`
	Status               OrderStatus `json:"status"`
}

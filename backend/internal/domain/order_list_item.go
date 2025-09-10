package domain

// OrderListItem is a projection for listing orders with related info
// It matches the requested SQL columns while keeping portability by formatting in Go
// created_at is returned as DD/MM/YYYY string
// Addresses are concatenated into single strings
// size_code comes from package_types
// Note: field tags map to JSON response keys

type OrderListItem struct {
	ID                     uint        `json:"id"`
	OrderNumber            string      `json:"order_number"`
	CreatedAt              string      `json:"created_at"`
	FullName               string      `json:"full_name"`
	OriginFullAddress      string      `json:"origin_full_address"`
	DestinationFullAddress string      `json:"destination_full_address"`
	ActualWeightKg         float64     `json:"actual_weight_kg"`
	SizeCode               PackageSize `json:"size_code"`
	Status                 OrderStatus `json:"status"`
}

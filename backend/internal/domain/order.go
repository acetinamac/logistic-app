package domain

import "time"

type OrderStatus string

const (
	StatusCreado      OrderStatus = "creado"
	StatusRecolectado OrderStatus = "recolectado"
	StatusEnEstacion  OrderStatus = "en_estacion"
	StatusEnRuta      OrderStatus = "en_ruta"
	StatusEntregado   OrderStatus = "entregado"
	StatusCancelado   OrderStatus = "cancelado"
)

type PackageSize string

const (
	SizeS PackageSize = "S"
	SizeM PackageSize = "M"
	SizeL PackageSize = "L"
)

type Coordinates struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Address struct {
	Country string `json:"country"`
	State   string `json:"state"`
	City    string `json:"city"`
	ZipCode string `json:"zipcode"`
	Street  string `json:"street"`
	ExtNum  string `json:"ext_num"`
	IntNum  string `json:"int_num"`
	Remarks string `json:"remarks"`
}

type Order struct {
	ID               uint        `json:"id" gorm:"primaryKey"`
	CustomerID       uint        `json:"customer_id"`
	OriginCoord      Coordinates `json:"origin_coord" gorm:"embedded;embeddedPrefix:origin_"`
	DestinationCoord Coordinates `json:"destination_coord" gorm:"embedded;embeddedPrefix:dest_"`
	OriginAddr       Address     `json:"origin_address" gorm:"embedded;embeddedPrefix:origin_"`
	DestinationAddr  Address     `json:"destination_address" gorm:"embedded;embeddedPrefix:dest_"`
	ItemsCount       int         `json:"items_count"`
	WeightKg         float64     `json:"weight_kg"`
	Size             PackageSize `json:"size"`
	Status           OrderStatus `json:"status"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
}

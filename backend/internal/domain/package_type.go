package domain

import "time"

// Package types table
type PackageType struct {
	ID          uint        `json:"id" gorm:"primaryKey"`
	SizeCode    PackageSize `json:"size_code" gorm:"type:package_size_enum;unique;not null"`
	MaxWeightKg float64     `json:"max_weight_kg" gorm:"type:decimal(5,2);not null"`
	Description string      `json:"description" gorm:"type:text"`
	IsActive    bool        `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time   `json:"created_at"`
}

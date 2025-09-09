package domain

import "time"

type Role string

const (
	RoleClient Role = "client"
	RoleAdmin  Role = "admin"
)

type User struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	Email            string    `json:"email" gorm:"uniqueIndex;size:255"`
	Password         string    `json:"-"`
	Phone            string    `json:"phone" gorm:"size:20"`
	FullName         string    `json:"full_name" gorm:"size:255"`
	Role             Role      `json:"role" gorm:"type:user_role_enum;default:client;not null"`
	IsActive         bool      `json:"is_active" gorm:"default:true;not null"`
	DefaultAddressID *uint     `json:"default_address_id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	UpdatedBy        *uint     `json:"updated_by"`
}

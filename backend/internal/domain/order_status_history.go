package domain

import "time"

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

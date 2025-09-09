package repository

import (
	"errors"
	"logistics-app/backend/internal/domain"
	"logistics-app/backend/internal/infra/db"
	"time"

	"gorm.io/gorm"
)

type OrderGormRepo struct{ db *gorm.DB }

func NewOrderGormRepo(database *db.Database) *OrderGormRepo { return &OrderGormRepo{db: database.DB} }

func (r *OrderGormRepo) Create(o *domain.Order) error { return r.db.Create(o).Error }

func (r *OrderGormRepo) FindByID(id uint) (*domain.Order, error) {
	var o domain.Order
	if err := r.db.First(&o, id).Error; err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *OrderGormRepo) FindByCustomer(customerID uint) ([]domain.Order, error) {
	var list []domain.Order
	if err := r.db.Where("customer_id = ?", customerID).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *OrderGormRepo) FindAll() ([]domain.Order, error) {
	var list []domain.Order
	if err := r.db.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *OrderGormRepo) UpdateStatus(id uint, status domain.OrderStatus, changedBy uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var o domain.Order
		if err := tx.First(&o, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("order not found")
			}
			return err
		}
		prev := o.Status
		if err := tx.Model(&domain.Order{}).Where("id = ?", id).Updates(map[string]interface{}{"status": status, "updated_by": changedBy}).Error; err != nil {
			return err
		}
		h := domain.OrderStatusHistory{
			OrderID:        id,
			PreviousStatus: prev,
			NewStatus:      status,
			ChangedAt:      time.Now(),
			ChangedBy:      changedBy,
		}
		if err := tx.Create(&h).Error; err != nil {
			return err
		}
		return nil
	})
}

package repository

import (
	"errors"
	"logistics-app/backend/internal/domain"
	"logistics-app/backend/internal/infra/db"

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

func (r *OrderGormRepo) UpdateStatus(id uint, status domain.OrderStatus) error {
	res := r.db.Model(&domain.Order{}).Where("id = ?", id).Update("status", status)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("order not found")
	}
	return nil
}

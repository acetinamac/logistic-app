package repository

import (
	"errors"
	"logistics-app/backend/internal/domain"
	"logistics-app/backend/internal/infra/db"
	"strings"
	"time"

	"gorm.io/gorm"
)

type OrderGormRepo struct{ db *gorm.DB }

// internal struct for scanning joined rows
type orderJoinedRow struct {
	Id             uint
	OrderNumber    string
	CreatedAt      time.Time
	FullName       string
	AOStreet       string
	AOExterior     string
	AONeighborhood string
	AOCity         string
	AOPostal       string
	ADStreet       string
	ADExterior     string
	ADNeighborhood string
	ADCity         string
	ADPostal       string
	Quantity       uint
	ActualWeightKg float64
	SizeCode       domain.PackageSize
	Status         domain.OrderStatus
}

func NewOrderGormRepo(database *db.Database) *OrderGormRepo {
	return &OrderGormRepo{db: database.DB}
}

func (r *OrderGormRepo) Create(o *domain.Order) error {
	return r.db.Create(o).Error
}

func (r *OrderGormRepo) FindByID(id uint) (*domain.Order, error) {
	var o domain.Order

	if err := r.db.First(&o, id).Error; err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *OrderGormRepo) FindDetailByID(id uint) (*domain.OrderDetail, error) {
	var d domain.OrderDetail

	q := r.db.Table("orders as o").
		Select("o.id, o.order_number, o.created_at, u.id as user_id, u.full_name, o.origin_address_id, ao.street as ao_street, ao.exterior_number as ao_exterior, ao.neighborhood as ao_neighborhood, ao.city as ao_city, ao.postal_code as ao_postal, o.destination_address_id, ad.street as ad_street, ad.exterior_number as ad_exterior, ad.neighborhood as ad_neighborhood, ad.city as ad_city, ad.postal_code as ad_postal, o.quantity, o.actual_weight_kg, o.package_type_id, pt.size_code, o.observations, o.internal_notes, o.updated_at, o.status").
		Joins("inner join users u on o.customer_id = u.id").
		Joins("inner join addresses ao on o.origin_address_id = ao.id").
		Joins("inner join addresses ad on o.destination_address_id = ad.id").
		Joins("inner join package_types pt on o.package_type_id = pt.id").
		Where("o.id = ?", id)

	if err := q.Take(&d).Error; err != nil {
		return nil, err
	}

	return &d, nil
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

func (r *OrderGormRepo) findJoined(base *gorm.DB) ([]domain.OrderListItem, error) {
	var rows []orderJoinedRow
	q := base.Table("orders as o").
		Select("o.id, o.order_number, o.created_at, u.full_name, ao.street as ao_street, ao.exterior_number as ao_exterior, ao.neighborhood as ao_neighborhood, ao.city as ao_city, ao.postal_code as ao_postal, ad.street as ad_street, ad.exterior_number as ad_exterior, ad.neighborhood as ad_neighborhood, ad.city as ad_city, ad.postal_code as ad_postal, o.quantity, o.actual_weight_kg, pt.size_code, o.status").
		Joins("inner join users u on o.customer_id = u.id").
		Joins("inner join addresses ao on o.origin_address_id = ao.id").
		Joins("inner join addresses ad on o.destination_address_id = ad.id").
		Joins("inner join package_types pt on o.package_type_id = pt.id")
	if err := q.Scan(&rows).Error; err != nil {
		return nil, err
	}

	// map rows to DTO
	items := make([]domain.OrderListItem, 0, len(rows))
	for _, rrow := range rows {
		origin := strings.TrimSpace(strings.Join([]string{rrow.AOStreet, rrow.AOExterior, rrow.AONeighborhood, rrow.AOCity, rrow.AOPostal}, " "))
		dest := strings.TrimSpace(strings.Join([]string{rrow.ADStreet, rrow.ADExterior, rrow.ADNeighborhood, rrow.ADCity, rrow.ADPostal}, " "))
		item := domain.OrderListItem{
			ID:                     rrow.Id,
			OrderNumber:            rrow.OrderNumber,
			CreatedAt:              rrow.CreatedAt.Format("02/01/2006"),
			FullName:               rrow.FullName,
			OriginFullAddress:      origin,
			DestinationFullAddress: dest,
			Quantity:               rrow.Quantity,
			ActualWeightKg:         rrow.ActualWeightKg,
			SizeCode:               rrow.SizeCode,
			Status:                 rrow.Status,
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *OrderGormRepo) FindJoinedByCustomer(customerID uint) ([]domain.OrderListItem, error) {
	base := r.db.Where("o.customer_id = ?", customerID)
	return r.findJoined(base)
}

func (r *OrderGormRepo) FindJoinedAll() ([]domain.OrderListItem, error) {
	base := r.db
	return r.findJoined(base)
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

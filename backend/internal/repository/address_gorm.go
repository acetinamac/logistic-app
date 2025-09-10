package repository

import (
	"errors"
	"logistics-app/backend/internal/domain"
	"logistics-app/backend/internal/infra/db"

	"gorm.io/gorm"
)

type AddressGormRepo struct{ db *gorm.DB }

func NewAddressGormRepo(database *db.Database) *AddressGormRepo {
	return &AddressGormRepo{db: database.DB}
}

// Composite payloads
type AddressWithCoords struct {
	Address     domain.Address      `json:"address"`
	Coordinates *domain.Coordinates `json:"coordinates,omitempty"`
}

func (r *AddressGormRepo) CreateWithCoordinates(customerID uint, payload AddressWithCoords) (*domain.Address, *domain.Coordinates, error) {
	var createdAddr domain.Address
	var createdCoord *domain.Coordinates
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		var coordID *uint
		if payload.Coordinates != nil {
			c := *payload.Coordinates
			if err := tx.Create(&c).Error; err != nil {
				return err
			}
			createdCoord = &c
			coordID = &c.ID
		}
		a := payload.Address
		a.CustomerID = customerID
		a.CoordinateID = coordID
		if err := tx.Create(&a).Error; err != nil {
			return err
		}
		createdAddr = a
		return nil
	}); err != nil {
		return nil, nil, err
	}
	return &createdAddr, createdCoord, nil
}

func (r *AddressGormRepo) UpdateWithCoordinates(requesterID uint, isAdmin bool, id uint, payload AddressWithCoords) (*domain.Address, *domain.Coordinates, error) {
	var outAddr domain.Address
	var outCoord *domain.Coordinates
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		var existing domain.Address
		q := tx.Where("id = ?", id)
		if !isAdmin {
			q = q.Where("customer_id = ?", requesterID)
		}
		if err := q.First(&existing).Error; err != nil {
			return err
		}

		// Update or create coordinates if provided
		if payload.Coordinates != nil {
			if existing.CoordinateID != nil {
				// update existing coord
				var c domain.Coordinates
				if err := tx.First(&c, *existing.CoordinateID).Error; err != nil {
					return err
				}
				c.Latitude = payload.Coordinates.Latitude
				c.Longitude = payload.Coordinates.Longitude
				if err := tx.Save(&c).Error; err != nil {
					return err
				}
				outCoord = &c
			} else {
				c := *payload.Coordinates
				if err := tx.Create(&c).Error; err != nil {
					return err
				}
				existing.CoordinateID = &c.ID
				outCoord = &c
			}
		}

		// Update address fields
		u := map[string]interface{}{
			"street":          payload.Address.Street,
			"exterior_number": payload.Address.ExteriorNumber,
			"interior_number": payload.Address.InteriorNumber,
			"neighborhood":    payload.Address.Neighborhood,
			"postal_code":     payload.Address.PostalCode,
			"city":            payload.Address.City,
			"state":           payload.Address.State,
			"country":         payload.Address.Country,
		}
		if err := tx.Model(&domain.Address{}).Where("id = ?", existing.ID).Updates(u).Error; err != nil {
			return err
		}
		if err := tx.First(&outAddr, existing.ID).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, nil, err
	}
	return &outAddr, outCoord, nil
}

func (r *AddressGormRepo) FindByID(requesterID uint, isAdmin bool, id uint) (*domain.Address, error) {
	var a domain.Address
	q := r.db.Where("id = ?", id)
	if !isAdmin {
		q = q.Where("customer_id = ?", requesterID)
	}
	if err := q.First(&a).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AddressGormRepo) List(requesterID uint, isAdmin bool, includeInactive bool) ([]domain.Address, error) {
	var list []domain.Address
	q := r.db.Model(&domain.Address{})
	if !isAdmin {
		q = q.Where("customer_id = ?", requesterID).Where("is_active = ?", true)
	} else if !includeInactive {
		q = q.Where("is_active = ?", true)
	}
	if err := q.Order("id asc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *AddressGormRepo) ToggleActive(requesterID uint, isAdmin bool, id uint, active bool) error {
	// Only owner or admin can toggle
	q := r.db.Model(&domain.Address{}).Where("id = ?", id)
	if !isAdmin {
		q = q.Where("customer_id = ?", requesterID)
	}
	res := q.Update("is_active", active)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *AddressGormRepo) Delete(requesterID uint, isAdmin bool, id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var a domain.Address
		q := tx.Where("id = ?", id)
		if !isAdmin {
			q = q.Where("customer_id = ?", requesterID)
		}
		if err := q.First(&a).Error; err != nil {
			return err
		}
		// Ensure no orders reference this address
		var cnt int64
		if err := tx.Model(&domain.Order{}).Where("origin_address_id = ? OR destination_address_id = ?", id, id).Count(&cnt).Error; err != nil {
			return err
		}
		if cnt > 0 {
			return errors.New("address is referenced by orders and cannot be deleted")
		}
		// Delete associated coordinates if not referenced by other addresses (rare)
		if err := tx.Delete(&domain.Address{}, a.ID).Error; err != nil {
			return err
		}
		if a.CoordinateID != nil {
			var usage int64
			if err := tx.Model(&domain.Address{}).Where("coordinate_id = ?", *a.CoordinateID).Count(&usage).Error; err != nil {
				return err
			}
			if usage == 0 {
				_ = tx.Delete(&domain.Coordinates{}, *a.CoordinateID).Error
			}
		}
		return nil
	})
}

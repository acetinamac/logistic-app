package repository

import (
	"logistics-app/backend/internal/domain"
	"logistics-app/backend/internal/infra/db"

	"gorm.io/gorm"
)

type PackageTypeGormRepo struct{ db *gorm.DB }

func NewPackageTypeGormRepo(database *db.Database) *PackageTypeGormRepo {
	return &PackageTypeGormRepo{db: database.DB}
}

func (r *PackageTypeGormRepo) FindAll(includeInactive bool) ([]domain.PackageType, error) {
	var list []domain.PackageType
	q := r.db.Model(&domain.PackageType{})

	if !includeInactive {
		q = q.Where("is_active = ?", true)
	}

	if err := q.Order("id asc").Find(&list).Error; err != nil {
		return nil, err
	}

	return list, nil
}

func (r *PackageTypeGormRepo) SetActive(id uint, active bool) error {
	res := r.db.Model(&domain.PackageType{}).Where("id = ?", id).Update("is_active", active)
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

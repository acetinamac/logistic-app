package repository

import (
	"logistics-app/backend/internal/domain"
	"logistics-app/backend/internal/infra/db"
)

type UserGormRepo struct {
	db *db.Database
}

func NewUserGormRepo(database *db.Database) *UserGormRepo {
	return &UserGormRepo{db: database}
}

func (r *UserGormRepo) Create(u *domain.User) error {
	return r.db.Create(u).Error
}

func (r *UserGormRepo) FindByEmail(email string) (*domain.User, error) {
	var u domain.User
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserGormRepo) FindByID(id uint) (*domain.User, error) {
	var u domain.User
	// Only select allowed fields
	if err := r.db.Model(&domain.User{}).
		Select("id, email, role, phone, full_name, is_active, created_at, updated_at").
		First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserGormRepo) DeleteByID(id uint) error {

	return r.db.Delete(&domain.User{}, id).Error
}

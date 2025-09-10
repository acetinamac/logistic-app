package usecase

import (
	"errors"
	"logistics-app/backend/internal/domain"
)

type PackageTypeRepo interface {
	FindAll(includeInactive bool) ([]domain.PackageType, error)
	SetActive(id uint, active bool) error
}

type PackageTypeService struct{ repo PackageTypeRepo }

func NewPackageTypeService(r PackageTypeRepo) *PackageTypeService {
	return &PackageTypeService{repo: r}
}

func (s *PackageTypeService) List(includeInactive bool) ([]domain.PackageType, error) {
	return s.repo.FindAll(includeInactive)
}

func (s *PackageTypeService) ToggleActive(id uint, active bool) error {
	if id == 0 {
		return errors.New("id requerido")
	}
	return s.repo.SetActive(id, active)
}

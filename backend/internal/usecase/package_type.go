package usecase

import (
	"errors"
	"logistics-app/backend/internal/domain"
	"sync"
	"time"
)

type PackageTypeRepo interface {
	FindAll(includeInactive bool) ([]domain.PackageType, error)
	SetActive(id uint, active bool) error
}

type PackageTypeService struct {
	repo       PackageTypeRepo
	cache      map[uint]domain.PackageType
	mutex      sync.RWMutex
	lastUpdate time.Time
	cacheTTL   time.Duration
}

func NewPackageTypeService(r PackageTypeRepo) *PackageTypeService {
	return &PackageTypeService{
		repo:     r,
		cache:    make(map[uint]domain.PackageType),
		cacheTTL: 10 * time.Second,
	}
}

func (s *PackageTypeService) List(includeInactive bool) ([]domain.PackageType, error) {
	return s.repo.FindAll(includeInactive)
}

func (s *PackageTypeService) ToggleActive(id uint, active bool) error {
	if id == 0 {
		return errors.New("id requerido")
	}

	if err := s.repo.SetActive(id, active); err == nil {
		s.invalidateCache()
	}

	return s.repo.SetActive(id, active)
}

func (s *PackageTypeService) GetPackageTypes() (map[uint]domain.PackageType, error) {
	s.mutex.RLock()

	if len(s.cache) > 0 && time.Since(s.lastUpdate) < s.cacheTTL {
		defer s.mutex.RUnlock()
		return s.cache, nil
	}
	s.mutex.RUnlock()

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if len(s.cache) > 0 && time.Since(s.lastUpdate) < s.cacheTTL {
		return s.cache, nil
	}

	types, err := s.repo.FindAll(false)
	if err != nil {
		return nil, err
	}

	newCache := make(map[uint]domain.PackageType)
	for _, pt := range types {
		newCache[pt.ID] = pt
	}

	s.cache = newCache
	s.lastUpdate = time.Now()

	return s.cache, nil
}

func (s *PackageTypeService) ValidatePackageWeight(packageTypeID uint, weightKg float64) error {
	packageTypes, err := s.GetPackageTypes()
	if err != nil {
		return err
	}

	errMsgExceeds := "El peso del paquete excede el límite estándar de 25kg. Para envíos de este tipo, debe contactar a la empresa para generar un convenio especial"
	packageType, exists := packageTypes[packageTypeID]
	if !exists {
		if weightKg > 25 {
			return errors.New(errMsgExceeds)
		}
		return errors.New("Tipo de paquete no encontrado")
	}

	if !packageType.IsActive {
		return errors.New("Tipo de paquete no está activo")
	}

	if weightKg > packageType.MaxWeightKg {
		if weightKg > 25 {
			return errors.New(errMsgExceeds)
		}
		return errors.New("El peso del paquete excede el límite máximo para este tipo de paquete")
	}

	return nil
}

func (s *PackageTypeService) invalidateCache() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.cache = make(map[uint]domain.PackageType)
	s.lastUpdate = time.Now()
}

package usecase

import (
	"errors"
	"logistics-app/backend/internal/domain"
	"logistics-app/backend/internal/repository"
)

type AddressRepo interface {
	CreateWithCoordinates(customerID uint, payload repository.AddressWithCoords) (*domain.Address, *domain.Coordinates, error)
	UpdateWithCoordinates(requesterID uint, isAdmin bool, id uint, payload repository.AddressWithCoords) (*domain.Address, *domain.Coordinates, error)
	FindByID(requesterID uint, isAdmin bool, id uint) (*domain.Address, error)
	List(requesterID uint, isAdmin bool, includeInactive bool) ([]domain.Address, error)
	ToggleActive(requesterID uint, isAdmin bool, id uint, active bool) error
	Delete(requesterID uint, isAdmin bool, id uint) error
}

type AddressService struct{ repo AddressRepo }

func NewAddressService(r AddressRepo) *AddressService { return &AddressService{repo: r} }

// DTOs to accept combined JSON in handlers
type AddressRequest struct {
	Street         string `json:"street"`
	ExteriorNumber string `json:"exterior_number"`
	InteriorNumber string `json:"interior_number"`
	Neighborhood   string `json:"neighborhood"`
	PostalCode     string `json:"postal_code"`
	City           string `json:"city"`
	State          string `json:"state"`
	Country        string `json:"country"`
	IsActive       *bool  `json:"is_active,omitempty"`
	Coordinates    *struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"coordinates"`
}

func (s *AddressService) toRepoPayload(req AddressRequest) repository.AddressWithCoords {
	addr := domain.Address{
		Street:         req.Street,
		ExteriorNumber: req.ExteriorNumber,
		InteriorNumber: req.InteriorNumber,
		Neighborhood:   req.Neighborhood,
		PostalCode:     req.PostalCode,
		City:           req.City,
		State:          req.State,
		Country:        req.Country,
	}
	var coords *domain.Coordinates
	if req.Coordinates != nil {
		coords = &domain.Coordinates{Latitude: req.Coordinates.Latitude, Longitude: req.Coordinates.Longitude}
	}
	return repository.AddressWithCoords{Address: addr, Coordinates: coords}
}

func (s *AddressService) Create(customerID uint, req AddressRequest) (*domain.Address, *domain.Coordinates, error) {
	if customerID == 0 {
		return nil, nil, errors.New("customerID requerido")
	}
	if req.Street == "" || req.City == "" || req.State == "" {
		return nil, nil, errors.New("street, city y state son requeridos")
	}
	payload := s.toRepoPayload(req)
	addr, coords, err := s.repo.CreateWithCoordinates(customerID, payload)
	if err != nil {
		return nil, nil, err
	}
	// allow overriding is_active on create if provided
	if req.IsActive != nil {
		_ = s.repo.ToggleActive(customerID, false, addr.ID, *req.IsActive)
		addr.IsActive = *req.IsActive
	}
	return addr, coords, nil
}

func (s *AddressService) Update(requesterID uint, isAdmin bool, id uint, req AddressRequest) (*domain.Address, *domain.Coordinates, error) {
	if id == 0 {
		return nil, nil, errors.New("id requerido")
	}
	payload := s.toRepoPayload(req)
	addr, coords, err := s.repo.UpdateWithCoordinates(requesterID, isAdmin, id, payload)
	if err != nil {
		return nil, nil, err
	}
	if req.IsActive != nil {
		if err := s.repo.ToggleActive(requesterID, isAdmin, id, *req.IsActive); err != nil {
			return nil, nil, err
		}
		addr.IsActive = *req.IsActive
	}
	return addr, coords, nil
}

func (s *AddressService) Get(requesterID uint, isAdmin bool, id uint) (*domain.Address, error) {
	return s.repo.FindByID(requesterID, isAdmin, id)
}

func (s *AddressService) List(requesterID uint, role domain.Role, includeInactive bool, all bool) ([]domain.Address, error) {
	isAdmin := role == domain.RoleAdmin
	// For clients, force includeInactive=false and own records only
	return s.repo.List(requesterID, isAdmin && all, includeInactive && isAdmin)
}

func (s *AddressService) ToggleActive(requesterID uint, isAdmin bool, id uint, active bool) error {
	return s.repo.ToggleActive(requesterID, isAdmin, id, active)
}

func (s *AddressService) Delete(requesterID uint, isAdmin bool, id uint) error {
	return s.repo.Delete(requesterID, isAdmin, id)
}

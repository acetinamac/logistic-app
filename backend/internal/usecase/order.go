package usecase

import (
	"errors"
	"logistics-app/backend/internal/domain"
)

type OrderRepo interface {
	Create(o *domain.Order) error
	FindByID(id uint) (*domain.Order, error)
	FindByCustomer(customerID uint) ([]domain.Order, error)
	FindAll() ([]domain.Order, error)
	UpdateStatus(id uint, status domain.OrderStatus) error
}

type OrderService struct{ repo OrderRepo }

func NewOrderService(r OrderRepo) *OrderService { return &OrderService{repo: r} }

func (s *OrderService) FindAll() ([]domain.Order, error) { return s.repo.FindAll() }
func (s *OrderService) FindByCustomer(customerID uint) ([]domain.Order, error) {
	return s.repo.FindByCustomer(customerID)
}

func (s *OrderService) DetermineSize(weight float64) (domain.PackageSize, error) {
	if weight <= 5 {
		return domain.SizeS, nil
	} else if weight <= 15 {
		return domain.SizeM, nil
	} else if weight <= 25 {
		return domain.SizeL, nil
	}
	return "", errors.New("peso superior a 25 kg: contactar a la empresa para convenio especial")
}

func (s *OrderService) Create(o *domain.Order) error {
	// Validar coordenadas
	if o.OriginCoord.Lat < -90 || o.OriginCoord.Lat > 90 || o.OriginCoord.Lng < -180 || o.OriginCoord.Lng > 180 {
		return errors.New("coordenadas de origen inválidas")
	}
	if o.DestinationCoord.Lat < -90 || o.DestinationCoord.Lat > 90 || o.DestinationCoord.Lng < -180 || o.DestinationCoord.Lng > 180 {
		return errors.New("coordenadas de destino inválidas")
	}
	if o.ItemsCount <= 0 {
		return errors.New("items_count debe ser > 0")
	}
	if o.WeightKg <= 0 {
		return errors.New("weight_kg debe ser > 0")
	}
	// Determinar tamaño
	size, err := s.DetermineSize(o.WeightKg)
	if err != nil {
		return err
	}
	o.Size = size
	o.Status = domain.StatusCreado
	return s.repo.Create(o)
}

func (s *OrderService) UpdateStatus(id uint, status domain.OrderStatus) error {
	// Se puede agregar validación de transición si se requiere.
	return s.repo.UpdateStatus(id, status)
}

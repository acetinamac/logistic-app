package usecase

import (
	"errors"
	"fmt"
	"logistics-app/backend/internal/domain"
	"time"
)

type OrderRepo interface {
	Create(o *domain.Order) error
	FindByID(id uint) (*domain.Order, error)
	FindByCustomer(customerID uint) ([]domain.Order, error)
	FindAll() ([]domain.Order, error)
	UpdateStatus(id uint, status domain.OrderStatus, changedBy uint) error
	FindJoinedByCustomer(customerID uint) ([]domain.OrderListItem, error)
	FindJoinedAll() ([]domain.OrderListItem, error)
	FindDetailByID(id uint) (*domain.OrderDetail, error)
}

type PackageTypeValidator interface {
	ValidatePackageWeight(packageTypeID uint, weightKg float64) error
}

type OrderService struct {
	repo             OrderRepo
	packageValidator PackageTypeValidator
}

func NewOrderService(r OrderRepo, pv PackageTypeValidator) *OrderService {
	return &OrderService{
		repo:             r,
		packageValidator: pv,
	}
}

func (s *OrderService) FindAll() ([]domain.Order, error) {
	return s.repo.FindAll()
}

func (s *OrderService) FindByCustomer(customerID uint) ([]domain.Order, error) {
	return s.repo.FindByCustomer(customerID)
}

func (s *OrderService) ListJoinedAll() ([]domain.OrderListItem, error) {

	return s.repo.FindJoinedAll()
}

func (s *OrderService) ListJoinedByCustomer(customerID uint) ([]domain.OrderListItem, error) {
	return s.repo.FindJoinedByCustomer(customerID)
}

func (s *OrderService) GetDetailByID(id uint) (*domain.OrderDetail, error) {
	return s.repo.FindDetailByID(id)
}

func generateOrderNumber(t time.Time) string {
	return fmt.Sprintf("ORD-%s-%d", t.Format("20060102"), t.UnixNano()%1_000_000)
}

func (s *OrderService) Create(o *domain.Order) error {
	if o.Quantity <= 0 {
		return errors.New("Quantity es requerido y debe ser mayor a 0")
	}

	if o.ActualWeightKg <= 0 {
		return errors.New("actual_weight_kg es requerido y debe ser mayor a 0")
	}

	if s.packageValidator != nil {
		if err := s.packageValidator.ValidatePackageWeight(o.PackageTypeID, o.ActualWeightKg); err != nil {
			return fmt.Errorf("Validación de peso: %w", err)
		}
	}

	if o.OriginAddressID == 0 || o.DestinationAddressID == 0 {
		return errors.New("origin_address_id y destination_address_id son requeridos")
	}

	if o.OriginAddressID == o.DestinationAddressID {
		return errors.New("Origin y destination deben ser diferentes")
	}

	if o.PackageTypeID == 0 {
		return errors.New("package_type_id es requerido")
	}

	if o.CustomerID == 0 || o.CreatedBy == 0 {
		return errors.New("customer_id y created_by son requeridos")
	}

	if o.OrderNumber == "" {
		o.OrderNumber = generateOrderNumber(time.Now())
	}

	if o.Status == "" {
		o.Status = domain.OrderCreated
	}
	return s.repo.Create(o)
}

func (s *OrderService) UpdateStatus(id uint, status domain.OrderStatus, changedBy uint) error {
	if changedBy == 0 {
		return errors.New("changedBy requerido")
	}

	return s.repo.UpdateStatus(id, status, changedBy)
}

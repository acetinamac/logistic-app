package tests

import (
	"errors"
	"logistics-app/backend/internal/domain"
	"logistics-app/backend/internal/usecase"
	"testing"
	"time"
)

type mockOrderRepo struct {
	orders     []domain.Order
	shouldFail bool
	failError  error
}

type mockPackageTypeValidator struct {
	shouldFail   bool
	failError    error
	packageTypes map[uint]domain.PackageType
}

func (m *mockPackageTypeValidator) ValidatePackageWeight(packageTypeID uint, weightKg float64) error {
	if m.shouldFail {
		return m.failError
	}

	if m.packageTypes != nil {
		packageType, exists := m.packageTypes[packageTypeID]
		if !exists {
			return errors.New("tipo de paquete no encontrado")
		}

		if !packageType.IsActive {
			return errors.New("tipo de paquete no está activo")
		}

		if weightKg > 25 {
			return errors.New("el peso del paquete excede el límite estándar de 25kg. Para envíos de este tipo, debe contactar a la empresa para generar un convenio especial")
		}

		if weightKg > packageType.MaxWeightKg {
			return errors.New("el peso del paquete excede el límite máximo para este tipo de paquete")
		}
	}

	return nil
}

func (m *mockOrderRepo) FindByID(id uint) (*domain.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockOrderRepo) FindByCustomer(customerID uint) ([]domain.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockOrderRepo) FindAll() ([]domain.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockOrderRepo) UpdateStatus(id uint, status domain.OrderStatus, changedBy uint) error {
	//TODO implement me
	panic("implement me")
}

func (m *mockOrderRepo) FindJoinedByCustomer(customerID uint) ([]domain.OrderListItem, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockOrderRepo) FindJoinedAll() ([]domain.OrderListItem, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockOrderRepo) FindDetailByID(id uint) (*domain.OrderDetail, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockOrderRepo) Create(o *domain.Order) error {
	if m.shouldFail {
		return m.failError
	}

	o.ID = uint(len(m.orders) + 1)
	o.CreatedAt = time.Now()
	o.UpdatedAt = time.Now()
	m.orders = append(m.orders, *o)

	return nil
}

func TestOrderService_Create_Success(t *testing.T) {
	// Arrange
	mockRepo := &mockOrderRepo{}
	mockValidator := &mockPackageTypeValidator{
		packageTypes: map[uint]domain.PackageType{
			1: {
				ID:          1,
				SizeCode:    domain.PackageM,
				MaxWeightKg: 5.0,
				IsActive:    true,
			},
		},
	}
	service := usecase.NewOrderService(mockRepo, mockValidator)

	order := &domain.Order{
		OriginAddressID:      1,
		DestinationAddressID: 2,
		PackageTypeID:        1,
		CustomerID:           1,
		CreatedBy:            1,
		Quantity:             1,
		ActualWeightKg:       2.5,
		Observations:         "Test order",
	}

	// Act
	err := service.Create(order)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if order.OrderNumber == "" {
		t.Error("Expected OrderNumber to be generated")
	}

	if order.Status != domain.OrderCreated {
		t.Errorf("Expected status to be %v, got %v", domain.OrderCreated, order.Status)
	}

	if len(mockRepo.orders) != 1 {
		t.Errorf("Expected 1 order in repository, got %d", len(mockRepo.orders))
	}
}

func TestOrderService_Create_MissingWeight(t *testing.T) {
	// Arrange
	mockRepo := &mockOrderRepo{}
	mockValidator := &mockPackageTypeValidator{}
	service := usecase.NewOrderService(mockRepo, mockValidator)

	order := &domain.Order{
		OriginAddressID:      1,
		DestinationAddressID: 2,
		PackageTypeID:        1,
		CustomerID:           1,
		CreatedBy:            1,
		Quantity:             1,
		ActualWeightKg:       0,
	}

	// Act
	err := service.Create(order)

	// Assert
	if err == nil {
		t.Error("Expected error for missing weight, got nil")
	}

	expectedError := "actual_weight_kg es requerido y debe ser mayor a 0"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestOrderService_Create_NegativeWeight(t *testing.T) {
	// Arrange
	mockRepo := &mockOrderRepo{}
	mockValidator := &mockPackageTypeValidator{}
	service := usecase.NewOrderService(mockRepo, mockValidator)

	order := &domain.Order{
		OriginAddressID:      1,
		DestinationAddressID: 2,
		PackageTypeID:        1,
		CustomerID:           1,
		CreatedBy:            1,
		Quantity:             1,
		ActualWeightKg:       -1.5,
	}

	// Act
	err := service.Create(order)

	// Assert
	if err == nil {
		t.Error("Expected error for negative weight, got nil")
	}

	expectedError := "actual_weight_kg es requerido y debe ser mayor a 0"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestOrderService_Create_PackageTypeNotFound(t *testing.T) {
	// Arrange
	mockRepo := &mockOrderRepo{}
	mockValidator := &mockPackageTypeValidator{
		shouldFail: true,
		failError:  errors.New("tipo de paquete no encontrado"),
	}
	service := usecase.NewOrderService(mockRepo, mockValidator)

	order := &domain.Order{
		OriginAddressID:      1,
		DestinationAddressID: 2,
		PackageTypeID:        999,
		CustomerID:           1,
		CreatedBy:            1,
		Quantity:             1,
		ActualWeightKg:       2.5,
	}

	// Act
	err := service.Create(order)

	// Assert
	if err == nil {
		t.Error("Expected error for package type not found, got nil")
	}

	if !errors.Is(err, errors.New("tipo de paquete no encontrado")) &&
		err.Error() != "Validación de peso: tipo de paquete no encontrado" {
		t.Errorf("Expected package type not found error, got '%s'", err.Error())
	}
}

func TestOrderService_Create_WeightExceedsPackageLimit(t *testing.T) {
	// Arrange
	mockRepo := &mockOrderRepo{}
	mockValidator := &mockPackageTypeValidator{
		packageTypes: map[uint]domain.PackageType{
			1: {
				ID:          1,
				SizeCode:    domain.PackageS,
				MaxWeightKg: 2.0,
				IsActive:    true,
			},
		},
	}
	service := usecase.NewOrderService(mockRepo, mockValidator)

	order := &domain.Order{
		OriginAddressID:      1,
		DestinationAddressID: 2,
		PackageTypeID:        1,
		CustomerID:           1,
		CreatedBy:            1,
		Quantity:             1,
		ActualWeightKg:       5.0,
	}

	// Act
	err := service.Create(order)

	// Assert
	if err == nil {
		t.Error("Expected error for weight exceeding package limit, got nil")
	}

	expectedError := "Validación de peso: el peso del paquete excede el límite máximo para este tipo de paquete"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestOrderService_Create_WeightExceedsStandardLimit(t *testing.T) {
	// Arrange
	mockRepo := &mockOrderRepo{}
	mockValidator := &mockPackageTypeValidator{
		packageTypes: map[uint]domain.PackageType{
			1: {
				ID:          1,
				SizeCode:    domain.PackageXL,
				MaxWeightKg: 30.0,
				IsActive:    true,
			},
		},
	}
	service := usecase.NewOrderService(mockRepo, mockValidator)

	order := &domain.Order{
		OriginAddressID:      1,
		DestinationAddressID: 2,
		PackageTypeID:        1,
		CustomerID:           1,
		CreatedBy:            1,
		Quantity:             1,
		ActualWeightKg:       28.0,
	}

	// Act
	err := service.Create(order)

	// Assert
	if err == nil {
		t.Error("Expected error for weight exceeding standard limit, got nil")
	}

	expectedError := "Validación de peso: el peso del paquete excede el límite estándar de 25kg. Para envíos de este tipo, debe contactar a la empresa para generar un convenio especial"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestOrderService_Create_MissingOriginAddress(t *testing.T) {
	// Arrange
	mockRepo := &mockOrderRepo{}
	mockValidator := &mockPackageTypeValidator{}
	service := usecase.NewOrderService(mockRepo, mockValidator)

	order := &domain.Order{
		DestinationAddressID: 2,
		PackageTypeID:        1,
		CustomerID:           1,
		CreatedBy:            1,
		Quantity:             1,
		ActualWeightKg:       2.5,
	}

	// Act
	err := service.Create(order)

	// Assert
	if err == nil {
		t.Error("Expected error for missing origin address, got nil")
	}

	expectedError := "origin_address_id y destination_address_id son requeridos"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestOrderService_Create_MissingDestinationAddress(t *testing.T) {
	// Arrange
	mockRepo := &mockOrderRepo{}
	mockValidator := &mockPackageTypeValidator{}
	service := usecase.NewOrderService(mockRepo, mockValidator)

	order := &domain.Order{
		OriginAddressID: 1,
		PackageTypeID:   1,
		CustomerID:      1,
		CreatedBy:       1,
		Quantity:        1,
		ActualWeightKg:  2.5,
	}

	// Act
	err := service.Create(order)

	// Assert
	if err == nil {
		t.Error("Expected error for missing destination address, got nil")
	}

	expectedError := "origin_address_id y destination_address_id son requeridos"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestOrderService_Create_SameOriginAndDestination(t *testing.T) {
	// Arrange
	mockRepo := &mockOrderRepo{}
	mockValidator := &mockPackageTypeValidator{}
	service := usecase.NewOrderService(mockRepo, mockValidator)

	order := &domain.Order{
		OriginAddressID:      1,
		DestinationAddressID: 1,
		PackageTypeID:        1,
		CustomerID:           1,
		CreatedBy:            1,
		Quantity:             1,
		ActualWeightKg:       2.5,
	}

	// Act
	err := service.Create(order)

	// Assert
	if err == nil {
		t.Error("Expected error for same origin and destination, got nil")
	}

	expectedError := "Origin y destination deben ser diferentes"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestOrderService_Create_MissingPackageType(t *testing.T) {
	// Arrange
	mockRepo := &mockOrderRepo{}
	mockValidator := &mockPackageTypeValidator{}
	service := usecase.NewOrderService(mockRepo, mockValidator)

	order := &domain.Order{
		OriginAddressID:      1,
		DestinationAddressID: 2,
		CustomerID:           1,
		CreatedBy:            1,
		Quantity:             1,
		ActualWeightKg:       2.5,
	}

	// Act
	err := service.Create(order)

	// Assert
	if err == nil {
		t.Error("Expected error for missing package type, got nil")
	}

	expectedError := "package_type_id es requerido"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestOrderService_Create_MissingCustomerID(t *testing.T) {
	// Arrange
	mockRepo := &mockOrderRepo{}
	mockValidator := &mockPackageTypeValidator{}
	service := usecase.NewOrderService(mockRepo, mockValidator)

	order := &domain.Order{
		OriginAddressID:      1,
		DestinationAddressID: 2,
		PackageTypeID:        1,
		CreatedBy:            1,
		Quantity:             1,
		ActualWeightKg:       2.5,
	}

	// Act
	err := service.Create(order)

	// Assert
	if err == nil {
		t.Error("Expected error for missing customer ID, got nil")
	}

	expectedError := "customer_id y created_by son requeridos"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestOrderService_Create_WithCustomOrderNumber(t *testing.T) {
	// Arrange
	mockRepo := &mockOrderRepo{}
	mockValidator := &mockPackageTypeValidator{}
	service := usecase.NewOrderService(mockRepo, mockValidator)

	customOrderNumber := "CUSTOM-12345"
	order := &domain.Order{
		OrderNumber:          customOrderNumber,
		OriginAddressID:      1,
		DestinationAddressID: 2,
		PackageTypeID:        1,
		CustomerID:           1,
		CreatedBy:            1,
		Quantity:             1,
		ActualWeightKg:       2.5,
	}

	// Act
	err := service.Create(order)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if order.OrderNumber != customOrderNumber {
		t.Errorf("Expected custom order number '%s', got '%s'", customOrderNumber, order.OrderNumber)
	}
}

func TestOrderService_Create_WithCustomStatus(t *testing.T) {
	// Arrange
	mockRepo := &mockOrderRepo{}
	mockValidator := &mockPackageTypeValidator{}
	service := usecase.NewOrderService(mockRepo, mockValidator)

	customStatus := domain.OrderCollected
	order := &domain.Order{
		OriginAddressID:      1,
		DestinationAddressID: 2,
		PackageTypeID:        1,
		CustomerID:           1,
		CreatedBy:            1,
		Status:               customStatus,
		Quantity:             1,
		ActualWeightKg:       2.5,
	}

	// Act
	err := service.Create(order)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if order.Status != customStatus {
		t.Errorf("Expected custom status '%s', got '%s'", customStatus, order.Status)
	}
}

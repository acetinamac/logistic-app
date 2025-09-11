package tests

import (
	"errors"
	"logistics-app/backend/internal/domain"
	"logistics-app/backend/internal/repository"
	"logistics-app/backend/internal/usecase"
	"testing"
	"time"
)

type mockAddressRepo struct {
	addresses   []domain.Address
	coordinates []domain.Coordinates
	shouldFail  bool
	failError   error
	lastAddrID  uint
	lastCoordID uint
}

func (m *mockAddressRepo) CreateWithCoordinates(customerID uint, payload repository.AddressWithCoords) (*domain.Address, *domain.Coordinates, error) {
	if m.shouldFail {
		return nil, nil, m.failError
	}

	// Create address
	m.lastAddrID++
	addr := payload.Address
	addr.ID = m.lastAddrID
	addr.CustomerID = customerID
	addr.CreatedAt = time.Now()
	addr.UpdatedAt = time.Now()
	addr.IsActive = true

	var coords *domain.Coordinates
	if payload.Coordinates != nil {
		m.lastCoordID++
		coord := *payload.Coordinates
		coord.ID = m.lastCoordID
		coord.CreatedAt = time.Now()
		coords = &coord
		addr.CoordinateID = &coord.ID
		m.coordinates = append(m.coordinates, coord)
	}

	m.addresses = append(m.addresses, addr)
	return &addr, coords, nil
}

func (m *mockAddressRepo) UpdateWithCoordinates(requesterID uint, isAdmin bool, id uint, payload repository.AddressWithCoords) (*domain.Address, *domain.Coordinates, error) {
	// Mock implementation for completeness
	return nil, nil, errors.New("not implemented in mock")
}

func (m *mockAddressRepo) FindByID(requesterID uint, isAdmin bool, id uint) (*domain.Address, error) {
	// Mock implementation for completeness
	return nil, errors.New("not implemented in mock")
}

func (m *mockAddressRepo) List(requesterID uint, isAdmin bool, includeInactive bool) ([]domain.Address, error) {
	// Mock implementation for completeness
	return nil, errors.New("not implemented in mock")
}

func (m *mockAddressRepo) ToggleActive(requesterID uint, isAdmin bool, id uint, active bool) error {
	for i, addr := range m.addresses {
		if addr.ID == id && (isAdmin || addr.CustomerID == requesterID) {
			m.addresses[i].IsActive = active
			return nil
		}
	}
	return errors.New("address not found")
}

func (m *mockAddressRepo) Delete(requesterID uint, isAdmin bool, id uint) error {
	// Mock implementation for completeness
	return errors.New("not implemented in mock")
}

func TestAddressService_Create_Success(t *testing.T) {
	// Arrange
	mockRepo := &mockAddressRepo{}
	service := usecase.NewAddressService(mockRepo)

	customerID := uint(1)
	req := usecase.AddressRequest{
		Street:         "Av. Reforma 123",
		ExteriorNumber: "123",
		InteriorNumber: "A",
		Neighborhood:   "Centro",
		PostalCode:     "06000",
		City:           "Ciudad de México",
		State:          "CDMX",
		Country:        "México",
		Coordinates: &struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		}{
			Latitude:  19.4326,
			Longitude: -99.1332,
		},
	}

	// Act
	addr, coords, err := service.Create(customerID, req)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if addr == nil {
		t.Fatal("Expected address to be created, got nil")
	}

	if coords == nil {
		t.Fatal("Expected coordinates to be created, got nil")
	}

	if addr.Street != req.Street {
		t.Errorf("Expected street '%s', got '%s'", req.Street, addr.Street)
	}

	if addr.City != req.City {
		t.Errorf("Expected city '%s', got '%s'", req.City, addr.City)
	}

	if addr.State != req.State {
		t.Errorf("Expected state '%s', got '%s'", req.State, addr.State)
	}

	if addr.CustomerID != customerID {
		t.Errorf("Expected customer ID %d, got %d", customerID, addr.CustomerID)
	}

	if coords.Latitude != req.Coordinates.Latitude {
		t.Errorf("Expected latitude %f, got %f", req.Coordinates.Latitude, coords.Latitude)
	}

	if coords.Longitude != req.Coordinates.Longitude {
		t.Errorf("Expected longitude %f, got %f", req.Coordinates.Longitude, coords.Longitude)
	}

	if len(mockRepo.addresses) != 1 {
		t.Errorf("Expected 1 address in repository, got %d", len(mockRepo.addresses))
	}
}

func TestAddressService_Create_WithoutCoordinates(t *testing.T) {
	// Arrange
	mockRepo := &mockAddressRepo{}
	service := usecase.NewAddressService(mockRepo)

	customerID := uint(1)
	req := usecase.AddressRequest{
		Street:  "Calle Principal 456",
		City:    "Guadalajara",
		State:   "Jalisco",
		Country: "México",
	}

	// Act
	addr, coords, err := service.Create(customerID, req)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if addr == nil {
		t.Fatal("Expected address to be created, got nil")
	}

	if coords != nil {
		t.Error("Expected coordinates to be nil, got non-nil")
	}

	if len(mockRepo.addresses) != 1 {
		t.Errorf("Expected 1 address in repository, got %d", len(mockRepo.addresses))
	}

	if len(mockRepo.coordinates) != 0 {
		t.Errorf("Expected 0 coordinates in repository, got %d", len(mockRepo.coordinates))
	}
}

func TestAddressService_Create_MissingCustomerID(t *testing.T) {
	// Arrange
	mockRepo := &mockAddressRepo{}
	service := usecase.NewAddressService(mockRepo)

	req := usecase.AddressRequest{
		Street: "Test Street",
		City:   "Test City",
		State:  "Test State",
	}

	// Act
	addr, coords, err := service.Create(0, req)

	// Assert
	if err == nil {
		t.Error("Expected error for missing customer ID, got nil")
	}

	expectedError := "customerID requerido"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}

	if addr != nil || coords != nil {
		t.Error("Expected both address and coordinates to be nil")
	}

	if len(mockRepo.addresses) != 0 {
		t.Errorf("Expected 0 addresses in repository, got %d", len(mockRepo.addresses))
	}
}

func TestAddressService_Create_MissingStreet(t *testing.T) {
	// Arrange
	mockRepo := &mockAddressRepo{}
	service := usecase.NewAddressService(mockRepo)

	customerID := uint(1)
	req := usecase.AddressRequest{
		City:  "Test City",
		State: "Test State",
	}

	// Act
	addr, coords, err := service.Create(customerID, req)

	// Assert
	if err == nil {
		t.Error("Expected error for missing street, got nil")
	}

	expectedError := "street, city y state son requeridos"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}

	if addr != nil || coords != nil {
		t.Error("Expected both address and coordinates to be nil")
	}
}

func TestAddressService_Create_MissingCity(t *testing.T) {
	// Arrange
	mockRepo := &mockAddressRepo{}
	service := usecase.NewAddressService(mockRepo)

	customerID := uint(1)
	req := usecase.AddressRequest{
		Street: "Test Street",
		State:  "Test State",
	}

	// Act
	addr, coords, err := service.Create(customerID, req)

	// Assert
	if err == nil {
		t.Error("Expected error for missing city, got nil")
	}

	expectedError := "street, city y state son requeridos"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}

	if addr != nil || coords != nil {
		t.Error("Expected both address and coordinates to be nil")
	}
}

func TestAddressService_Create_MissingState(t *testing.T) {
	// Arrange
	mockRepo := &mockAddressRepo{}
	service := usecase.NewAddressService(mockRepo)

	customerID := uint(1)
	req := usecase.AddressRequest{
		Street: "Test Street",
		City:   "Test City",
	}

	// Act
	addr, coords, err := service.Create(customerID, req)

	// Assert
	if err == nil {
		t.Error("Expected error for missing state, got nil")
	}

	expectedError := "street, city y state son requeridos"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}

	if addr != nil || coords != nil {
		t.Error("Expected both address and coordinates to be nil")
	}
}

func TestAddressService_Create_InvalidLatitudeMin(t *testing.T) {
	// Arrange
	mockRepo := &mockAddressRepo{}
	service := usecase.NewAddressService(mockRepo)

	customerID := uint(1)
	req := usecase.AddressRequest{
		Street: "Test Street",
		City:   "Test City",
		State:  "Test State",
		Coordinates: &struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		}{
			Latitude:  -91.0, // Invalid: less than -90
			Longitude: 0.0,
		},
	}

	// Act
	addr, coords, err := service.Create(customerID, req)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid latitude, got nil")
	}

	expectedError := "latitud debe estar entre -90 y 90 grados"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}

	if addr != nil || coords != nil {
		t.Error("Expected both address and coordinates to be nil")
	}
}

func TestAddressService_Create_InvalidLatitudeMax(t *testing.T) {
	// Arrange
	mockRepo := &mockAddressRepo{}
	service := usecase.NewAddressService(mockRepo)

	customerID := uint(1)
	req := usecase.AddressRequest{
		Street: "Test Street",
		City:   "Test City",
		State:  "Test State",
		Coordinates: &struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		}{
			Latitude:  91.0, // Invalid: greater than 90
			Longitude: 0.0,
		},
	}

	// Act
	addr, coords, err := service.Create(customerID, req)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid latitude, got nil")
	}

	expectedError := "latitud debe estar entre -90 y 90 grados"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}

	if addr != nil || coords != nil {
		t.Error("Expected both address and coordinates to be nil")
	}
}

func TestAddressService_Create_InvalidLongitudeMin(t *testing.T) {
	// Arrange
	mockRepo := &mockAddressRepo{}
	service := usecase.NewAddressService(mockRepo)

	customerID := uint(1)
	req := usecase.AddressRequest{
		Street: "Test Street",
		City:   "Test City",
		State:  "Test State",
		Coordinates: &struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		}{
			Latitude:  0.0,
			Longitude: -181.0, // Invalid: less than -180
		},
	}

	// Act
	addr, coords, err := service.Create(customerID, req)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid longitude, got nil")
	}

	expectedError := "longitud debe estar entre -180 y 180 grados"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}

	if addr != nil || coords != nil {
		t.Error("Expected both address and coordinates to be nil")
	}
}

func TestAddressService_Create_InvalidLongitudeMax(t *testing.T) {
	// Arrange
	mockRepo := &mockAddressRepo{}
	service := usecase.NewAddressService(mockRepo)

	customerID := uint(1)
	req := usecase.AddressRequest{
		Street: "Test Street",
		City:   "Test City",
		State:  "Test State",
		Coordinates: &struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		}{
			Latitude:  0.0,
			Longitude: 181.0, // Invalid: greater than 180
		},
	}

	// Act
	addr, coords, err := service.Create(customerID, req)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid longitude, got nil")
	}

	expectedError := "longitud debe estar entre -180 y 180 grados"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}

	if addr != nil || coords != nil {
		t.Error("Expected both address and coordinates to be nil")
	}
}

func TestAddressService_Create_ValidCoordinatesAtBoundaries(t *testing.T) {
	// Arrange
	mockRepo := &mockAddressRepo{}
	service := usecase.NewAddressService(mockRepo)

	testCases := []struct {
		name      string
		latitude  float64
		longitude float64
	}{
		{"North Pole", 90.0, 0.0},
		{"South Pole", -90.0, 0.0},
		{"Prime Meridian", 0.0, 0.0},
		{"International Date Line East", 0.0, 180.0},
		{"International Date Line West", 0.0, -180.0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			customerID := uint(1)
			req := usecase.AddressRequest{
				Street: "Test Street",
				City:   "Test City",
				State:  "Test State",
				Coordinates: &struct {
					Latitude  float64 `json:"latitude"`
					Longitude float64 `json:"longitude"`
				}{
					Latitude:  tc.latitude,
					Longitude: tc.longitude,
				},
			}

			// Act
			addr, coords, err := service.Create(customerID, req)

			// Assert
			if err != nil {
				t.Errorf("Expected no error for valid coordinates, got %v", err)
			}

			if addr == nil {
				t.Fatal("Expected address to be created, got nil")
			}

			if coords == nil {
				t.Fatal("Expected coordinates to be created, got nil")
			}

			if coords.Latitude != tc.latitude {
				t.Errorf("Expected latitude %f, got %f", tc.latitude, coords.Latitude)
			}

			if coords.Longitude != tc.longitude {
				t.Errorf("Expected longitude %f, got %f", tc.longitude, coords.Longitude)
			}
		})
	}
}

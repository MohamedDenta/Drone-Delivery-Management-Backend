package service

import (
	"github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/domain"
	"github.com/stretchr/testify/mock"
)

// MockDroneRepository is a mock of DroneRepository
type MockDroneRepository struct {
	mock.Mock
}

func (m *MockDroneRepository) CreateDrone(drone *domain.Drone) error {
	args := m.Called(drone)
	return args.Error(0)
}

func (m *MockDroneRepository) GetActiveDrones() ([]*domain.Drone, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Drone), args.Error(1)
}

func (m *MockDroneRepository) GetIdleDrones() ([]*domain.Drone, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Drone), args.Error(1)
}

func (m *MockDroneRepository) GetDroneByID(id string) (*domain.Drone, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Drone), args.Error(1)
}

func (m *MockDroneRepository) GetDroneByName(name string) (*domain.Drone, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Drone), args.Error(1)
}

func (m *MockDroneRepository) UpdateDrone(drone *domain.Drone) error {
	args := m.Called(drone)
	return args.Error(0)
}

// MockOrderRepository is a mock of OrderRepository
type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) CreateOrder(order *domain.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockOrderRepository) GetOrderByID(id string) (*domain.Order, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Order), args.Error(1)
}

func (m *MockOrderRepository) GetNextPendingOrder() (*domain.Order, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Order), args.Error(1)
}

func (m *MockOrderRepository) ClaimNextPendingOrder(droneID string) (*domain.Order, error) {
	args := m.Called(droneID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Order), args.Error(1)
}

func (m *MockOrderRepository) GetActiveOrderByDroneID(droneID string) (*domain.Order, error) {
	args := m.Called(droneID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Order), args.Error(1)
}

func (m *MockOrderRepository) UpdateOrder(order *domain.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

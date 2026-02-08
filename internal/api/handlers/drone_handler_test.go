package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/domain"
	"github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDroneRepo struct {
	mock.Mock
}

func (m *MockDroneRepo) CreateDrone(drone *domain.Drone) error {
	args := m.Called(drone)
	return args.Error(0)
}
func (m *MockDroneRepo) GetDroneByID(id string) (*domain.Drone, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Drone), args.Error(1)
}
func (m *MockDroneRepo) GetDroneByName(name string) (*domain.Drone, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Drone), args.Error(1)
}
func (m *MockDroneRepo) GetIdleDrones() ([]*domain.Drone, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Drone), args.Error(1)
}
func (m *MockDroneRepo) GetActiveDrones() ([]*domain.Drone, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Drone), args.Error(1)
}
func (m *MockDroneRepo) GetAllDrones() ([]*domain.Drone, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Drone), args.Error(1)
}
func (m *MockDroneRepo) UpdateDrone(drone *domain.Drone) error {
	args := m.Called(drone)
	return args.Error(0)
}

func TestRegister_Endpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(MockDroneRepo)
	droneService := service.NewDroneService(mockRepo, nil)
	handler := NewDroneHandler(droneService, nil)

	r := gin.New()
	r.POST("/drones", handler.Register)

	reqBody := RegisterDroneRequest{Name: "Test-Drone"}
	body, _ := json.Marshal(reqBody)

	mockRepo.On("GetDroneByName", "Test-Drone").Return(nil, domain.ErrNotFound)
	mockRepo.On("CreateDrone", mock.Anything).Return(nil)

	req, _ := http.NewRequest(http.MethodPost, "/drones", bytes.NewBuffer(body))
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
}

func TestListDrones_Endpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(MockDroneRepo)
	droneService := service.NewDroneService(mockRepo, nil)
	handler := NewDroneHandler(droneService, nil)

	r := gin.New()
	r.GET("/drones", handler.ListDrones)

	mockRepo.On("GetAllDrones").Return([]*domain.Drone{{Name: "D1"}}, nil)

	req, _ := http.NewRequest(http.MethodGet, "/drones", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

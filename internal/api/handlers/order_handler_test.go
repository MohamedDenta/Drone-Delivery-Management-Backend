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
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Reuse mocks from service package (need to be exported or redefined)
// Since they are in the 'service' package and these tests are in 'handlers',
// and 'mocks_test.go' is likely not exported, I'll redefine what I need or
// use a common mock package. For simplicity, I'll redefine a minimal mock here
// or use the real service with a mock repo.

type MockOrderRepo struct {
	mock.Mock
}

func (m *MockOrderRepo) CreateOrder(order *domain.Order) error {
	args := m.Called(order)
	return args.Error(0)
}
func (m *MockOrderRepo) GetOrderByID(id string) (*domain.Order, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Order), args.Error(1)
}
func (m *MockOrderRepo) GetActiveOrderByDroneID(droneID string) (*domain.Order, error) {
	args := m.Called(droneID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Order), args.Error(1)
}
func (m *MockOrderRepo) GetNextPendingOrder() (*domain.Order, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Order), args.Error(1)
}
func (m *MockOrderRepo) ClaimNextPendingOrder(droneID string) (*domain.Order, error) {
	args := m.Called(droneID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Order), args.Error(1)
}
func (m *MockOrderRepo) GetAllOrders() ([]*domain.Order, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Order), args.Error(1)
}
func (m *MockOrderRepo) UpdateOrder(order *domain.Order) error {
	args := m.Called(order)
	return args.Error(0)
}
func (m *MockOrderRepo) UpdateOrderCoords(id string, oLat, oLon, dLat, dLon float64) error {
	args := m.Called(id, oLat, oLon, dLat, dLon)
	return args.Error(0)
}

func TestCreateOrder_Endpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(MockOrderRepo)
	orderService := service.NewOrderService(mockRepo, nil)
	handler := NewOrderHandler(orderService)

	r := gin.New()
	r.POST("/orders", handler.CreateOrder)

	reqBody := CreateOrderRequest{
		OriginLat: 10, OriginLon: 10, DestLat: 20, DestLon: 20,
	}
	body, _ := json.Marshal(reqBody)

	mockRepo.On("CreateOrder", mock.Anything).Return(nil)

	req, _ := http.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
	var order domain.Order
	json.Unmarshal(resp.Body.Bytes(), &order)
	assert.Equal(t, domain.OrderStatusPending, order.Status)
}

func TestGetOrder_Endpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(MockOrderRepo)
	orderService := service.NewOrderService(mockRepo, nil)
	handler := NewOrderHandler(orderService)

	r := gin.New()
	r.GET("/orders/:id", handler.GetOrder)

	orderID := ksuid.New()
	mockRepo.On("GetOrderByID", orderID.String()).Return(&domain.Order{ID: orderID, Status: domain.OrderStatusPending}, nil)

	req, _ := http.NewRequest(http.MethodGet, "/orders/"+orderID.String(), nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockRepo.AssertExpectations(t)
}

package service

import (
	"testing"
	"time"

	"github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/domain"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateOrder(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	service := NewOrderService(mockRepo)

	mockRepo.On("CreateOrder", mock.AnythingOfType("*domain.Order")).Return(nil)

	order, err := service.CreateOrder(1.0, 1.0, 2.0, 2.0)

	assert.NoError(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, domain.OrderStatusPending, order.Status)
	assert.NotEqual(t, ksuid.Nil, order.ID)
	mockRepo.AssertExpectations(t)
}

func TestUpdateOrderState_ValidTransition(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	service := NewOrderService(mockRepo)

	orderID := ksuid.New()
	existingOrder := &domain.Order{
		ID:        orderID,
		Status:    domain.OrderStatusPending,
		UpdatedAt: time.Now(),
	}

	mockRepo.On("GetOrderByID", orderID.String()).Return(existingOrder, nil)
	mockRepo.On("UpdateOrder", mock.MatchedBy(func(o *domain.Order) bool {
		return o.Status == domain.OrderStatusReserved && o.ID == orderID
	})).Return(nil)

	updatedOrder, err := service.UpdateOrderState(orderID.String(), domain.OrderStatusReserved)

	assert.NoError(t, err)
	assert.Equal(t, domain.OrderStatusReserved, updatedOrder.Status)
	mockRepo.AssertExpectations(t)
}

func TestUpdateOrderState_InvalidTransition(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	service := NewOrderService(mockRepo)

	orderID := ksuid.New()
	existingOrder := &domain.Order{
		ID:     orderID,
		Status: domain.OrderStatusPending,
	}

	mockRepo.On("GetOrderByID", orderID.String()).Return(existingOrder, nil)

	// Trying to go straight to DELIVERED from PENDING should fail
	_, err := service.UpdateOrderState(orderID.String(), domain.OrderStatusDelivered)

	assert.Error(t, err)
	assert.Equal(t, "invalid state transition", err.Error())
	mockRepo.AssertNotCalled(t, "UpdateOrder")
}

package service

import (
	"errors"
	"testing"

	"github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/domain"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/mock"
)

func TestRecoveryHandler_OnDroneStatusChanged_Success(t *testing.T) {
	mockOrderRepo := new(MockOrderRepository)
	handler := NewRecoveryHandler(mockOrderRepo)

	droneID := ksuid.New()
	orderID := ksuid.New()
	lat, lon := 50.0, 10.0

	existingOrder := &domain.Order{
		ID:      orderID,
		Status:  domain.OrderStatusPickedUp,
		DroneID: &droneID,
	}

	// Expect searching for active order
	mockOrderRepo.On("GetActiveOrderByDroneID", droneID.String()).Return(existingOrder, nil)

	// Expect updating order
	mockOrderRepo.On("UpdateOrder", mock.MatchedBy(func(o *domain.Order) bool {
		return o.ID == orderID &&
			o.Status == domain.OrderStatusPending &&
			o.DroneID == nil &&
			o.OriginLat == lat && o.OriginLon == lon
	})).Return(nil)

	handler.OnDroneStatusChanged(droneID.String(), domain.DroneStatusDelivering, domain.DroneStatusBroken, lat, lon)

	mockOrderRepo.AssertExpectations(t)
}

func TestRecoveryHandler_OnDroneStatusChanged_Ignored(t *testing.T) {
	mockOrderRepo := new(MockOrderRepository)
	handler := NewRecoveryHandler(mockOrderRepo)

	// Not BROKEN
	handler.OnDroneStatusChanged("id", domain.DroneStatusIdle, domain.DroneStatusDelivering, 0, 0)

	// BROKEN but not from Delivering (e.g. from IDLE)
	handler.OnDroneStatusChanged("id", domain.DroneStatusIdle, domain.DroneStatusBroken, 0, 0)

	mockOrderRepo.AssertNotCalled(t, "GetActiveOrderByDroneID")
}

func TestRecoveryHandler_OnDroneStatusChanged_NoActiveOrder(t *testing.T) {
	mockOrderRepo := new(MockOrderRepository)
	handler := NewRecoveryHandler(mockOrderRepo)

	droneID := ksuid.New().String()

	// Return NotFound
	mockOrderRepo.On("GetActiveOrderByDroneID", droneID).Return(nil, domain.ErrNotFound)

	handler.OnDroneStatusChanged(droneID, domain.DroneStatusDelivering, domain.DroneStatusBroken, 0, 0)

	mockOrderRepo.AssertExpectations(t)
	mockOrderRepo.AssertNotCalled(t, "UpdateOrder")
}

func TestRecoveryHandler_OnDroneStatusChanged_RepoError(t *testing.T) {
	mockOrderRepo := new(MockOrderRepository)
	handler := NewRecoveryHandler(mockOrderRepo)

	droneID := ksuid.New().String()

	// Return Generic Error
	mockOrderRepo.On("GetActiveOrderByDroneID", droneID).Return(nil, errors.New("db error"))

	handler.OnDroneStatusChanged(droneID, domain.DroneStatusDelivering, domain.DroneStatusBroken, 0, 0)

	mockOrderRepo.AssertExpectations(t)
	mockOrderRepo.AssertNotCalled(t, "UpdateOrder")
}

package service

import (
	"testing"
	"time"

	"github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/domain"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestReserveJob_Success(t *testing.T) {
	mockDroneRepo := new(MockDroneRepository)
	mockOrderRepo := new(MockOrderRepository)
	dispatcher := NewDispatcherService(mockDroneRepo, mockOrderRepo)

	droneID := ksuid.New()
	drone := &domain.Drone{
		ID:     droneID,
		Status: domain.DroneStatusIdle,
	}

	orderID := ksuid.New()
	order := &domain.Order{
		ID:        orderID,
		Status:    domain.OrderStatusPending,
		UpdatedAt: time.Now(),
	}

	mockDroneRepo.On("GetDroneByID", droneID.String()).Return(drone, nil)
	mockOrderRepo.On("GetNextPendingOrder").Return(order, nil)

	// Expect Order Update
	mockOrderRepo.On("UpdateOrder", mock.MatchedBy(func(o *domain.Order) bool {
		return o.Status == domain.OrderStatusReserved && *o.DroneID == droneID
	})).Return(nil)

	// Expect Drone Update
	mockDroneRepo.On("UpdateDrone", mock.MatchedBy(func(d *domain.Drone) bool {
		return d.Status == domain.DroneStatusDelivering
	})).Return(nil)

	assignedOrder, err := dispatcher.ReserveJob(droneID.String())

	assert.NoError(t, err)
	assert.NotNil(t, assignedOrder)
	assert.Equal(t, domain.OrderStatusReserved, assignedOrder.Status)
	assert.Equal(t, droneID, *assignedOrder.DroneID)

	mockDroneRepo.AssertExpectations(t)
	mockOrderRepo.AssertExpectations(t)
}

func TestReserveJob_NoDroneFound(t *testing.T) {
	mockDroneRepo := new(MockDroneRepository)
	mockOrderRepo := new(MockOrderRepository)
	dispatcher := NewDispatcherService(mockDroneRepo, mockOrderRepo)

	mockDroneRepo.On("GetDroneByID", "invalid-id").Return(nil, domain.ErrNotFound) // Assuming domain has ErrNotFound or repo returns generic error? Repo returns arbitrary error.
	// Actually repo defines ErrNotFound in repository package, but we mock it.
	// Let's return errors.New("not found") or just verify behavior on error.

	_, err := dispatcher.ReserveJob("invalid-id")
	assert.Error(t, err)
}

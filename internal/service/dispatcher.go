package service

import (
	"errors"
	"time"

	"github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/domain"
	"github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/repository"
)

type DispatcherService struct {
	droneRepo repository.DroneRepository
	orderRepo repository.OrderRepository
}

func NewDispatcherService(droneRepo repository.DroneRepository, orderRepo repository.OrderRepository) *DispatcherService {
	return &DispatcherService{
		droneRepo: droneRepo,
		orderRepo: orderRepo,
	}
}

// ReserveJob assigns the next pending order to the requesting drone
func (s *DispatcherService) ReserveJob(droneID string) (*domain.Order, error) {
	// 1. Get Drone
	drone, err := s.droneRepo.GetDroneByID(droneID)
	if err != nil {
		return nil, err
	}

	// 2. Validate Drone Status
	if drone.Status != domain.DroneStatusIdle {
		return nil, errors.New("drone is not idle")
	}

	// 3. Get Next Pending Order
	// TODO: In a real system, this should be a transaction to prevent race conditions.
	// For this implementation, we rely on the database's atomic operations where possible,
	// but a proper SELECT FOR UPDATE would be better.
	order, err := s.orderRepo.GetNextPendingOrder()
	if err != nil {
		return nil, errors.New("no pending orders available")
	}

	// 4. Assign Order (This part needs concurrency control in production)
	order.Status = domain.OrderStatusReserved
	order.DroneID = &drone.ID
	order.UpdatedAt = time.Now()

	if err := s.orderRepo.UpdateOrder(order); err != nil {
		return nil, err
	}

	// 5. Update Drone Status
	drone.Status = domain.DroneStatusDelivering
	if err := s.droneRepo.UpdateDrone(drone); err != nil {
		// Rollback order assignment if drone update fails (manual compensation)
		order.Status = domain.OrderStatusPending
		order.DroneID = nil
		_ = s.orderRepo.UpdateOrder(order)
		return nil, err
	}

	return order, nil
}

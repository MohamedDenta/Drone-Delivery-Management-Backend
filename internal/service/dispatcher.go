package service

import (
	"errors"

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

	// 3. Claim Next Pending Order (Atomic)
	order, err := s.orderRepo.ClaimNextPendingOrder(drone.ID.String())
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, errors.New("no pending orders available")
		}
		return nil, err
	}

	// 4. Update Drone Status
	drone.Status = domain.DroneStatusDelivering
	if err := s.droneRepo.UpdateDrone(drone); err != nil {
		// Rollback order assignment if drone update fails
		order.Status = domain.OrderStatusPending
		order.DroneID = nil
		_ = s.orderRepo.UpdateOrder(order)
		return nil, err
	}

	return order, nil
}

package service

import (
	"log"
	"time"

	"github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/domain"
	"github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/repository"
)

type RecoveryHandler struct {
	orderRepo repository.OrderRepository
}

func NewRecoveryHandler(orderRepo repository.OrderRepository) *RecoveryHandler {
	return &RecoveryHandler{orderRepo: orderRepo}
}

func (h *RecoveryHandler) OnDroneStatusChanged(droneID string, oldStatus, newStatus domain.DroneStatus, currentLat, currentLon float64) {
	// Broken Drone Recovery Logic
	if newStatus == domain.DroneStatusBroken && oldStatus == domain.DroneStatusDelivering {
		// 1. Find the active order
		order, err := h.orderRepo.GetActiveOrderByDroneID(droneID)
		if err == nil {
			// 2. Update Order: Origin becomes current drone location
			order.OriginLat = currentLat
			order.OriginLon = currentLon

			// 3. Reset Status to PENDING and unassign drone
			order.Status = domain.OrderStatusPending
			order.DroneID = nil
			order.UpdatedAt = time.Now()

			if err := h.orderRepo.UpdateOrder(order); err != nil {
				log.Printf("Failed to recover order for broken drone %s: %v", droneID, err)
			} else {
				log.Printf("Recovered order %s from broken drone %s", order.ID, droneID)
			}
		} else if err != domain.ErrNotFound {
			log.Printf("Error finding active order for broken drone %s: %v", droneID, err)
		}
	}
}

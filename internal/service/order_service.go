package service

import (
	"errors"
	"time"

	"github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/domain"
	"github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/repository"
	"github.com/segmentio/ksuid"
)

type OrderService struct {
	repo repository.OrderRepository
}

func NewOrderService(repo repository.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrder(originLat, originLon, destLat, destLon float64) (*domain.Order, error) {
	order := &domain.Order{
		ID:        ksuid.New(),
		Status:    domain.OrderStatusPending,
		OriginLat: originLat,
		OriginLon: originLon,
		DestLat:   destLat,
		DestLon:   destLon,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.CreateOrder(order); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrderService) GetOrder(id string) (*domain.Order, error) {
	return s.repo.GetOrderByID(id)
}

func (s *OrderService) UpdateOrderState(id string, newState domain.OrderStatus) (*domain.Order, error) {
	order, err := s.repo.GetOrderByID(id)
	if err != nil {
		return nil, err
	}

	// Validate state transition (Simple valid transitions)
	if !isValidTransition(order.Status, newState) {
		return nil, errors.New("invalid state transition")
	}

	order.Status = newState
	order.UpdatedAt = time.Now()

	if err := s.repo.UpdateOrder(order); err != nil {
		return nil, err
	}
	return order, nil
}

func isValidTransition(current, next domain.OrderStatus) bool {
	switch current {
	case domain.OrderStatusPending:
		return next == domain.OrderStatusReserved
	case domain.OrderStatusReserved:
		return next == domain.OrderStatusPickedUp
	case domain.OrderStatusPickedUp:
		return next == domain.OrderStatusDelivered || next == domain.OrderStatusFailed
	case domain.OrderStatusFailed, domain.OrderStatusDelivered:
		return false // Terminal states
	default:
		return false
	}
}

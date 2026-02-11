package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/domain"
	"github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/infrastructure/rabbitmq"
	"github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/repository"
	"github.com/segmentio/ksuid"
)

type OrderService struct {
	repo      repository.OrderRepository
	publisher *rabbitmq.Client
}

func NewOrderService(repo repository.OrderRepository, publisher *rabbitmq.Client) *OrderService {
	return &OrderService{
		repo:      repo,
		publisher: publisher,
	}
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

	// Publish OrderCreated event
	if s.publisher != nil {
		event := domain.OrderCreatedEvent{
			OrderID:   order.ID.String(),
			OriginLat: order.OriginLat,
			OriginLon: order.OriginLon,
			DestLat:   order.DestLat,
			DestLon:   order.DestLon,
			Timestamp: time.Now(),
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.publisher.Publish(ctx, "order.created", event); err != nil {
			log.Printf("Failed to publish OrderCreated event for order %s: %v", order.ID.String(), err)
		}
	}

	return order, nil
}

func (s *OrderService) GetOrder(id string) (*domain.Order, error) {
	order, err := s.repo.GetOrderByID(id)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrderService) ListOrders() ([]*domain.Order, error) {
	return s.repo.GetAllOrders()
}

func (s *OrderService) WithdrawOrder(id string) error {
	order, err := s.repo.GetOrderByID(id)
	if err != nil {
		return err
	}

	if order.Status != domain.OrderStatusPending && order.Status != domain.OrderStatusReserved {
		return errors.New("cannot withdraw order that is already picked up or finished")
	}

	order.Status = domain.OrderStatusCancelled
	order.UpdatedAt = time.Now()
	return s.repo.UpdateOrder(order)
}

func (s *OrderService) UpdateOrderCoords(id string, originLat, originLon, destLat, destLon float64) error {
	order, err := s.repo.GetOrderByID(id)
	if err != nil {
		return err
	}

	if order.Status != domain.OrderStatusPending {
		return errors.New("cannot update destination of an order that is already in progress")
	}

	return s.repo.UpdateOrderCoords(id, originLat, originLon, destLat, destLon)
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
	if next == domain.OrderStatusCancelled {
		return current == domain.OrderStatusPending || current == domain.OrderStatusReserved
	}

	switch current {
	case domain.OrderStatusPending:
		return next == domain.OrderStatusReserved
	case domain.OrderStatusReserved:
		return next == domain.OrderStatusPickedUp
	case domain.OrderStatusPickedUp:
		return next == domain.OrderStatusDelivered || next == domain.OrderStatusFailed
	case domain.OrderStatusFailed, domain.OrderStatusDelivered, domain.OrderStatusCancelled:
		return false // Terminal states
	default:
		return false
	}
}

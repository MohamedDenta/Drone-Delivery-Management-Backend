package service

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/domain"
	"github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/infrastructure/rabbitmq"
	"github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/repository"
)

type OrderDispatcherWorker struct {
	rabbitClient *rabbitmq.Client
	dispatcher   *DispatcherService
	droneRepo    repository.DroneRepository
}

func NewOrderDispatcherWorker(rabbitClient *rabbitmq.Client, dispatcher *DispatcherService, droneRepo repository.DroneRepository) *OrderDispatcherWorker {
	return &OrderDispatcherWorker{
		rabbitClient: rabbitClient,
		dispatcher:   dispatcher,
		droneRepo:    droneRepo,
	}
}

func (w *OrderDispatcherWorker) Start() error {
	if w.rabbitClient == nil {
		return nil
	}

	return w.rabbitClient.Subscribe("order_dispatch_queue", "order.created", func(body []byte) error {
		var event domain.OrderCreatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		log.Printf("Worker received OrderCreated event for ID: %s. Attempting to find a drone...", event.OrderID)

		// 1. Find Idle Drones
		drones, err := w.droneRepo.GetIdleDrones()
		if err != nil {
			return err
		}

		if len(drones) == 0 {
			log.Printf("No idle drones available for order %s. Message will be requeued.", event.OrderID)
			// Return error to requeue message
			return errors.New("no idle drones available")
		}

		// 2. Simple assignment (First available).
		drone := drones[0]

		_, err = w.dispatcher.ReserveJob(drone.ID.String())
		if err != nil {
			log.Printf("Failed to reserve job for drone %s: %v", drone.ID.String(), err)
			return err
		}

		log.Printf("Successfully assigned order %s to drone %s", event.OrderID, drone.ID.String())
		return nil
	})
}

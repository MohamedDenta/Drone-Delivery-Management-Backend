package service

import (
	"context"
	"log"
	"time"

	"github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/domain"
	infra "github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/infrastructure/redis"
	"github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/repository"
)

type HeartbeatMonitor struct {
	droneRepo    repository.DroneRepository
	droneService *DroneService
	redisClient  *infra.Client
	interval     time.Duration
}

func NewHeartbeatMonitor(repo repository.DroneRepository, svc *DroneService, redis *infra.Client) *HeartbeatMonitor {
	return &HeartbeatMonitor{
		droneRepo:    repo,
		droneService: svc,
		redisClient:  redis,
		interval:     10 * time.Second, // Scan every 10 seconds
	}
}

func (m *HeartbeatMonitor) Start(ctx context.Context) {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	log.Println("Heartbeat Monitor started")

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.checkDrones()
		}
	}
}

func (m *HeartbeatMonitor) checkDrones() {
	drones, err := m.droneRepo.GetActiveDrones()
	if err != nil {
		log.Printf("HeartbeatMonitor: failed to get active drones: %v", err)
		return
	}

	for _, drone := range drones {
		if m.redisClient == nil {
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		alive, err := m.redisClient.HasDroneHeartbeat(ctx, drone.ID.String())
		cancel()

		if err != nil {
			log.Printf("HeartbeatMonitor: failed to check heartbeat for drone %s: %v", drone.ID.String(), err)
			continue
		}

		if !alive {
			log.Printf("HeartbeatMonitor: Drone %s is OFFLINE", drone.ID.String())
			// Update status to OFFLINE via DroneService (which triggers observers)
			if err := m.droneService.UpdateStatus(drone.ID.String(), domain.DroneStatusOffline); err != nil {
				log.Printf("HeartbeatMonitor: failed to update status for drone %s: %v", drone.ID.String(), err)
			}
		}
	}
}

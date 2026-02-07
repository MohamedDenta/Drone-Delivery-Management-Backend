package service

import (
	"errors"
	"time"

	"github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/domain"
	"github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/repository"
	"github.com/segmentio/ksuid"
)

// DroneStatusObserver defines the interface for listening to drone status changes
type DroneStatusObserver interface {
	OnDroneStatusChanged(droneID string, oldStatus, newStatus domain.DroneStatus, currentLat, currentLon float64)
}

type DroneService struct {
	repo      repository.DroneRepository
	observers []DroneStatusObserver
}

func NewDroneService(repo repository.DroneRepository) *DroneService {
	return &DroneService{
		repo:      repo,
		observers: make([]DroneStatusObserver, 0),
	}
}

func (s *DroneService) AddObserver(observer DroneStatusObserver) {
	s.observers = append(s.observers, observer)
}

func (s *DroneService) RegisterDrone(name string) (*domain.Drone, error) {
	// check if exists
	if _, err := s.repo.GetDroneByName(name); err == nil {
		return nil, errors.New("drone already exists")
	}

	drone := &domain.Drone{
		ID:        ksuid.New(),
		Name:      name,
		Status:    domain.DroneStatusIdle,
		Latitude:  0,
		Longitude: 0,
		CreatedAt: time.Now(),
	}

	if err := s.repo.CreateDrone(drone); err != nil {
		return nil, err
	}
	return drone, nil
}

func (s *DroneService) UpdateLocation(id string, lat, lon float64) error {
	drone, err := s.repo.GetDroneByID(id)
	if err != nil {
		return err
	}
	drone.Latitude = lat
	drone.Longitude = lon
	return s.repo.UpdateDrone(drone)
}

func (s *DroneService) UpdateStatus(id string, status domain.DroneStatus) error {
	drone, err := s.repo.GetDroneByID(id)
	if err != nil {
		return err
	}

	// Broken Drone Recovery Logic
	// Notify observers
	for _, observer := range s.observers {
		observer.OnDroneStatusChanged(id, drone.Status, status, drone.Latitude, drone.Longitude)
	}

	drone.Status = status
	return s.repo.UpdateDrone(drone)
}

func (s *DroneService) GetDrone(id string) (*domain.Drone, error) {
	return s.repo.GetDroneByID(id)
}

func (s *DroneService) GetDroneByName(name string) (*domain.Drone, error) {
	return s.repo.GetDroneByName(name)
}

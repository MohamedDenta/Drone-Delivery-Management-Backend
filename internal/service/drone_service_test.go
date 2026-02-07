package service

import (
	"errors"
	"testing"

	"github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/domain"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockObserver
type MockDroneStatusObserver struct {
	mock.Mock
}

func (m *MockDroneStatusObserver) OnDroneStatusChanged(droneID string, oldStatus, newStatus domain.DroneStatus, currentLat, currentLon float64) {
	m.Called(droneID, oldStatus, newStatus, currentLat, currentLon)
}

func TestRegisterDrone_Success(t *testing.T) {
	mockRepo := new(MockDroneRepository)
	service := NewDroneService(mockRepo)

	name := "Drone-01"

	// Expect GetDroneByName to return Error (Not Found), which means name is available
	mockRepo.On("GetDroneByName", name).Return(nil, errors.New("record not found"))

	// Expect CreateDrone to be called
	mockRepo.On("CreateDrone", mock.MatchedBy(func(d *domain.Drone) bool {
		return d.Name == name && d.Status == domain.DroneStatusIdle && d.ID != ksuid.Nil
	})).Return(nil)

	drone, err := service.RegisterDrone(name)

	assert.NoError(t, err)
	assert.NotNil(t, drone)
	assert.Equal(t, name, drone.Name)
	mockRepo.AssertExpectations(t)
}

func TestRegisterDrone_DuplicateName(t *testing.T) {
	mockRepo := new(MockDroneRepository)
	service := NewDroneService(mockRepo)

	name := "Drone-01"
	existingDrone := &domain.Drone{Name: name}

	// Expect GetDroneByName to return Success (Found), which means duplicate
	mockRepo.On("GetDroneByName", name).Return(existingDrone, nil)

	_, err := service.RegisterDrone(name)

	assert.Error(t, err)
	assert.Equal(t, "drone already exists", err.Error())
	mockRepo.AssertNotCalled(t, "CreateDrone")
}

func TestUpdateLocation(t *testing.T) {
	mockRepo := new(MockDroneRepository)
	service := NewDroneService(mockRepo)

	id := ksuid.New()
	existingDrone := &domain.Drone{ID: id, Latitude: 0, Longitude: 0}

	mockRepo.On("GetDroneByID", id.String()).Return(existingDrone, nil)
	mockRepo.On("UpdateDrone", mock.MatchedBy(func(d *domain.Drone) bool {
		return d.ID == id && d.Latitude == 10.5 && d.Longitude == 20.5
	})).Return(nil)

	err := service.UpdateLocation(id.String(), 10.5, 20.5)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateStatus_BrokenRescue_NotifyObserver(t *testing.T) {
	mockRepo := new(MockDroneRepository)
	mockObserver := new(MockDroneStatusObserver)
	service := NewDroneService(mockRepo)
	service.AddObserver(mockObserver)

	droneID := ksuid.New()
	existingDrone := &domain.Drone{ID: droneID, Status: domain.DroneStatusDelivering, Latitude: 50.0, Longitude: 10.0}

	mockRepo.On("GetDroneByID", droneID.String()).Return(existingDrone, nil)

	// Expect Drone Update
	mockRepo.On("UpdateDrone", mock.MatchedBy(func(d *domain.Drone) bool {
		return d.ID == droneID && d.Status == domain.DroneStatusBroken
	})).Return(nil)

	// Expect Observer Notification
	mockObserver.On("OnDroneStatusChanged", droneID.String(), domain.DroneStatusDelivering, domain.DroneStatusBroken, 50.0, 10.0).Return()

	err := service.UpdateStatus(droneID.String(), domain.DroneStatusBroken)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockObserver.AssertExpectations(t)
}

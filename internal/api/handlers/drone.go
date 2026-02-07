package handlers

import (
	"log/slog"
	"net/http"

	"github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/service"
	"github.com/gin-gonic/gin"
)

type DroneHandler struct {
	droneService      *service.DroneService
	dispatcherService *service.DispatcherService
}

func NewDroneHandler(droneService *service.DroneService, dispatcherService *service.DispatcherService) *DroneHandler {
	return &DroneHandler{
		droneService:      droneService,
		dispatcherService: dispatcherService,
	}
}

type RegisterDroneRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *DroneHandler) Register(c *gin.Context) {
	var req RegisterDroneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	drone, err := h.droneService.RegisterDrone(req.Name)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, drone)
}

// UpdateLocation handles heartbeat and location updates
func (h *DroneHandler) UpdateLocation(c *gin.Context) {

	name := c.MustGet("user").(string)

	var req struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Resolve Name -> Drone
	drone, err := h.droneService.GetDroneByName(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "drone not found"})
		return
	}

	// Update Location
	if err := h.droneService.UpdateLocation(drone.ID.String(), req.Latitude, req.Longitude); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "drone_id": drone.ID})
}

func (h *DroneHandler) ReserveJob(c *gin.Context) {

	type ReserveRequest struct {
		DroneID string `json:"drone_id" binding:"required"`
	}
	var req ReserveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.dispatcherService.ReserveJob(req.DroneID)
	if err != nil {
		slog.Error("failed to reserve job", "error", err)
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

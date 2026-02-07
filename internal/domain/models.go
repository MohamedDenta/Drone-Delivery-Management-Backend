package domain

import (
	"time"

	"github.com/segmentio/ksuid"
)

// DroneStatus represents the current state of a drone
type DroneStatus string

const (
	DroneStatusIdle       DroneStatus = "IDLE"
	DroneStatusDelivering DroneStatus = "DELIVERING"
	DroneStatusBroken     DroneStatus = "BROKEN"
)

// Drone represents a delivery drone in the system
type Drone struct {
	ID        ksuid.KSUID `json:"id"`
	Name      string      `json:"name"`
	Status    DroneStatus `json:"status"`
	Latitude  float64     `json:"latitude"`
	Longitude float64     `json:"longitude"`
	CreatedAt time.Time   `json:"created_at"`
}

// OrderStatus represents the lifecycle state of an order
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "PENDING"
	OrderStatusReserved  OrderStatus = "RESERVED"
	OrderStatusPickedUp  OrderStatus = "PICKED_UP"
	OrderStatusDelivered OrderStatus = "DELIVERED"
	OrderStatusFailed    OrderStatus = "FAILED"
)

// Order represents a delivery order
type Order struct {
	ID        ksuid.KSUID  `json:"id"`
	Status    OrderStatus  `json:"status"`
	OriginLat float64      `json:"origin_lat"`
	OriginLon float64      `json:"origin_lon"`
	DestLat   float64      `json:"dest_lat"`
	DestLon   float64      `json:"dest_lon"`
	DroneID   *ksuid.KSUID `json:"drone_id,omitempty"` // Nullable if not assigned
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

package domain

import "time"

type OrderCreatedEvent struct {
	OrderID   string    `json:"order_id"`
	OriginLat float64   `json:"origin_lat"`
	OriginLon float64   `json:"origin_lon"`
	DestLat   float64   `json:"dest_lat"`
	DestLon   float64   `json:"dest_lon"`
	Timestamp time.Time `json:"timestamp"`
}

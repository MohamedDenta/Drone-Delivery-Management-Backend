package api

import (
	"github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/api/handlers"
	"github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/api/middleware"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetupRouter(
	droneHandler *handlers.DroneHandler,
	orderHandler *handlers.OrderHandler,
) *gin.Engine {
	r := gin.New()

	// Middleware
	r.Use(gin.Recovery())
	r.Use(middleware.TelemetryMiddleware("drone-backend"))

	// Public Routes
	r.POST("/auth/token", handlers.Login)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Protected Routes
	api := r.Group("/api/v1")
	api.Use(middleware.AuthMiddleware())
	{
		// Drone Routes
		api.POST("/drones", droneHandler.Register)
		api.POST("/drones/location", droneHandler.UpdateLocation)
		api.POST("/drones/jobs/reserve", droneHandler.ReserveJob)

		// Order Routes
		api.POST("/orders", orderHandler.CreateOrder)
		api.POST("/orders/:id/status", orderHandler.UpdateStatus)
	}

	return r
}

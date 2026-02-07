package handlers

import (
	"log/slog"
	"net/http"

	"github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/domain"
	"github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/service"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

type CreateOrderRequest struct {
	OriginLat float64 `json:"origin_lat" binding:"required"`
	OriginLon float64 `json:"origin_lon" binding:"required"`
	DestLat   float64 `json:"dest_lat" binding:"required"`
	DestLon   float64 `json:"dest_lon" binding:"required"`
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.orderService.CreateOrder(req.OriginLat, req.OriginLon, req.DestLat, req.DestLon)
	if err != nil {
		slog.Error("failed to create order", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order"})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status domain.OrderStatus `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.orderService.UpdateOrderState(id, req.Status)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()}) // Conflict for invalid state transition
		return
	}

	c.JSON(http.StatusOK, order)
}
func (h *OrderHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")
	order, err := h.orderService.GetOrder(id)
	if err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, order)
}

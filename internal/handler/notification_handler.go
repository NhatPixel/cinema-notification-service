package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	
	"github.com/NhatPixel/cinema-notification-service/internal/service"
	"github.com/NhatPixel/cinema-notification-service/internal/model"
)

type NotificationHandler struct {
	service *service.NotificationService
}

func NewNotificationHandler(s *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: s}
}

func (h *NotificationHandler) Create(c *gin.Context) {
	var n model.Notification
	if err := c.ShouldBindJSON(&n); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Notify(n); err != nil {
		c.JSON(500, gin.H{"error": "failed to create notification"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "sent"})
}
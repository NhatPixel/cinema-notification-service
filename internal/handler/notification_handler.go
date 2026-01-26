package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	
	"github.com/NhatPixel/cinema-notification-service/internal/service"
	"github.com/NhatPixel/cinema-notification-service/internal/dto"
	"github.com/NhatPixel/cinema-notification-service/internal/validation"
)

type NotificationHandler struct {
	service *service.NotificationService
}

func NewNotificationHandler(s *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: s}
}

func (h *NotificationHandler) Create(c *gin.Context) {
	var req dto.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validation.TranslateValidationError(err).Error()})
		return
	}

	if err := h.service.Create(req); err != nil {
		c.JSON(500, gin.H{"error": "Tạo thông báo thất bại!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "sent"})
}

func (h *NotificationHandler) CreateForUsers(c *gin.Context) {
	var reqs []dto.CreateRequest
	if err := c.ShouldBindJSON(&reqs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validation.TranslateValidationError(err).Error()})
		return
	}

	if err := h.service.CreateForUsers(reqs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tạo thông báo thất bại!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Tạo thông báo thành công."})
}

func (h *NotificationHandler) UpdateReadStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Thiếu mã thông báo!"})
		return
	}

	if err := h.service.UpdateReadStatus(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Đọc thông báo thất bại!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "read"})
}

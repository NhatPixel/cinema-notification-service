package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/NhatPixel/cinema-notification-service/internal/service"
)

type SSEHandler struct {
	service *service.NotificationService
}

func NewSSEHandler(s *service.NotificationService) *SSEHandler {
	return &SSEHandler{service: s}
}

func (h *SSEHandler) Stream(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	ch := h.service.Subscribe(userID)
	defer h.service.Unsubscribe(userID, ch)

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get flusher"})
		return
	}

	for {
		select {
		case msg := <-ch:
			c.SSEvent("notification", msg)
			flusher.Flush()
		case <-c.Request.Context().Done():
			return
		}
	}
}
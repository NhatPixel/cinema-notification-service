package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/NhatPixel/cinema-notification-service/internal/dto"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Thiếu mã người dùng!"})
		return
	}

	notifications, err := h.service.FindByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tải thông báo thất bại!"})
		return
	}

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tải flusher thất bại!"})
		return
	}

	for _, n := range notifications {
		var resp dto.SSEResponse
		resp.FromModel(n)
		c.SSEvent("notification", resp)
		flusher.Flush()
	}

	ch := h.service.Subscribe(userID)
	defer h.service.Unsubscribe(userID, ch)

	for {
		select {
		case msg := <-ch:
			var resp dto.SSEResponse
			resp.FromModel(msg)
			c.SSEvent("notification", resp)
			flusher.Flush()
		case <-c.Request.Context().Done():
			return
		}
	}
}
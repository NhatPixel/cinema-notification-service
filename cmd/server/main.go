package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	
	"github.com/NhatPixel/cinema-notification-service/internal/handler"
	"github.com/NhatPixel/cinema-notification-service/internal/service"
	"github.com/NhatPixel/cinema-notification-service/config"
	"github.com/NhatPixel/cinema-notification-service/internal/repository"
)

func main() {
	godotenv.Load()
	
	r := gin.Default()

	db, err := config.NewMySQL()
	if err != nil {
		log.Fatal("cannot connect db:", err)
	}

	notificationRepo := repository.NewNotificationRepo(db)
	notificationService := service.NewNotificationService(notificationRepo)

	sseHandler := handler.NewSSEHandler(notificationService)
	notificationHandler := handler.NewNotificationHandler(notificationService)

	r.GET("/notifications/stream", sseHandler.Stream)
	r.POST("/notifications", notificationHandler.Create)
	r.POST("/notifications/bulk", notificationHandler.CreateForUsers)
	r.PUT("/notifications/:id/read", notificationHandler.UpdateReadStatus)	

	r.Run(":8080")
}
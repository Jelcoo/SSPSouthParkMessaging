package main

import (
	"os"
	"sp-messaging/api/internal/adapters/handler"
	"sp-messaging/api/internal/adapters/messaging"
	"sp-messaging/api/internal/core/services"

	"github.com/gin-gonic/gin"
)

func main() {
	rabbitURL := os.Getenv("RABBITMQ_URL")

	publisher := messaging.NewRabbitMQPublisher(rabbitURL, "shit-pipe")
	defer publisher.Close()

	svc := services.NewMessengerService(publisher)

	router := gin.Default()
	h := handler.NewHTTPHandler(svc)
	router.POST("/shit-post", h.PostMessage)

	port := os.Getenv("API_PORT")
	router.Run(":" + port)
}

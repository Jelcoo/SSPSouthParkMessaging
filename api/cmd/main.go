package main

import (
	"flag"
	"os"
	"sp-messaging/api/internal/adapters/handler"
	"sp-messaging/api/internal/adapters/messaging"
	"sp-messaging/api/internal/core/services"

	"github.com/gin-gonic/gin"
)

func main() {
	rabbitURL := flag.String("rabbitmq", envOr("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"), "RabbitMQ connection URL")
	flag.Parse()

	publisher := messaging.NewRabbitMQPublisher(*rabbitURL, "shit-pipe")
	defer publisher.Close()

	svc := services.NewMessengerService(publisher)

	router := gin.Default()
	h := handler.NewHTTPHandler(svc)
	router.POST("/shit-post", h.PostMessage)

	port := os.Getenv("API_PORT")
	router.Run(":" + port)
}

func envOr(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

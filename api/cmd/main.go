package main

import (
	"os"
	"sp-messaging/api/internal/adapters/handler"
	"sp-messaging/api/internal/adapters/messaging"
	"sp-messaging/api/internal/adapters/scheduler"
	"sp-messaging/api/internal/core/services"

	"github.com/gin-gonic/gin"
)

func main() {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	rabbitExchange := os.Getenv("RABBITMQ_EXCHANGE")

	publisher := messaging.NewRabbitMQPublisher(rabbitURL, rabbitExchange)
	defer publisher.Close()

	svc := services.NewMessengerService(publisher)

	cronSchedule := os.Getenv("CRON_SCHEDULE")
	cronScheduler := scheduler.NewCronScheduler(svc)
	if err := cronScheduler.ScheduleRandomMessages(cronSchedule); err != nil {
		panic(err)
	}
	cronScheduler.Start()
	defer cronScheduler.Stop()

	router := gin.Default()
	h := handler.NewHTTPHandler(svc)
	router.POST("/shit-post", h.PostMessage)

	port := os.Getenv("API_PORT")
	router.Run(":" + port)
}

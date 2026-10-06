package messaging

import (
	"context"
	"encoding/json"
	"sp-messaging/api/internal/core/domain"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQPublisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
}

func NewRabbitMQPublisher(url string, queue string) *RabbitMQPublisher {
	conn, err := amqp.Dial(url)
	if err != nil {
		panic(err)
	}

	channel, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	if _, err := channel.QueueDeclare(queue, true, false, false, false, nil); err != nil {
		panic(err)
	}

	return &RabbitMQPublisher{
		conn:    conn,
		channel: channel,
		queue:   queue,
	}
}

func (r *RabbitMQPublisher) PublishMessage(message domain.Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return r.channel.PublishWithContext(context.Background(), "", r.queue, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		MessageId:    message.ID,
		Body:         data,
	})
}

func (r *RabbitMQPublisher) Close() error {
	if err := r.channel.Close(); err != nil {
		return err
	}
	return r.conn.Close()
}

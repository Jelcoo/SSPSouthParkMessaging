package messaging

import (
	"context"
	"encoding/json"
	"regexp"
	"sp-messaging/api/internal/core/domain"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
)

const routingKeyPrefix = "messages"

var nonAlphanumeric = regexp.MustCompile(`[^a-z0-9]+`)

type RabbitMQPublisher struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
}

func NewRabbitMQPublisher(url string, exchange string) *RabbitMQPublisher {
	conn, err := amqp.Dial(url)
	if err != nil {
		panic(err)
	}

	channel, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	if err := channel.ExchangeDeclare(exchange, amqp.ExchangeTopic, true, false, false, false, nil); err != nil {
		panic(err)
	}

	return &RabbitMQPublisher{
		conn:     conn,
		channel:  channel,
		exchange: exchange,
	}
}

func (r *RabbitMQPublisher) PublishMessage(message domain.Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return r.channel.PublishWithContext(context.Background(), r.exchange, routingKey(message), false, false, amqp.Publishing{
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

func routingKey(message domain.Message) string {
	author := strings.Trim(nonAlphanumeric.ReplaceAllString(strings.ToLower(message.Author), "-"), "-")
	if author == "" {
		author = "unknown"
	}
	return routingKeyPrefix + "." + author
}

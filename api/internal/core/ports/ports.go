package ports

import "sp-messaging/api/internal/core/domain"

type MessengerService interface {
	PostMessage(message domain.Message) (*domain.Message, error)
}

type MessagePublisher interface {
	PublishMessage(message domain.Message) error
}

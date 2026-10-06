package services

import (
	"sp-messaging/api/internal/core/domain"
	"sp-messaging/api/internal/core/ports"

	"github.com/google/uuid"
)

type MessengerService struct {
	publisher ports.MessagePublisher
}

func NewMessengerService(publisher ports.MessagePublisher) *MessengerService {
	return &MessengerService{
		publisher: publisher,
	}
}

func (m *MessengerService) PostMessage(message domain.Message) (*domain.Message, error) {
	message.ID = uuid.New().String()
	if err := m.publisher.PublishMessage(message); err != nil {
		return nil, err
	}
	return &message, nil
}

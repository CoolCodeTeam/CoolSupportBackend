package useCase

import (
	chats "github.com/CoolCodeTeam/CoolSupportBackend/chats/usecase"
	"github.com/CoolCodeTeam/CoolSupportBackend/messages/models"
	"github.com/CoolCodeTeam/CoolSupportBackend/messages/repository"
)

//go:generate moq -out messages_ucase_mock.go . MessagesUseCase
type MessagesUseCase interface {
	SaveChatMessage(message *models.Message) (uint64, error)
	GetChatMessages(chatID uint64, userID uint64) (models.Messages, error)
}

type MessageUseCaseImpl struct {
	repository repository.MessageRepository
	chats      chats.ChatsUseCase
}

func NewMessageUseCase(repository repository.MessageRepository, chats chats.ChatsUseCase) MessagesUseCase {
	return &MessageUseCaseImpl{
		repository: repository,
		chats:      chats,
	}
}

func (m *MessageUseCaseImpl) GetChatMessages(chatID uint64, userID uint64) (models.Messages, error) {

	return m.repository.GetMessagesByChatID(chatID)
}




func (m *MessageUseCaseImpl) SaveChatMessage(message *models.Message) (uint64, error) {
	return m.repository.PutMessage(message)
}



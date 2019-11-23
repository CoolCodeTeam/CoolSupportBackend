package repository

import "github.com/CoolCodeTeam/CoolSupportBackend/messages/models"

//go:generate moq -out message_repo_mock.go . MessageRepository

type MessageRepository interface {
	PutMessage(message *models.Message) (uint64, error)
	GetMessagesByChatID(chatID uint64) (models.Messages, error)
}

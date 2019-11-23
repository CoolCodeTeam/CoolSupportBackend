package repository

import "github.com/CoolCodeTeam/CoolSupportBackend/notifications/models"

type NotificationRepository interface {
	GetNotificationHub(chatID uint64) *models.WebSocketHub
}

type NotificationArrayRepository struct {
	Hubs map[uint64]*models.WebSocketHub
}

func (n *NotificationArrayRepository) GetNotificationHub(chatID uint64) *models.WebSocketHub {
	if hub, ok := n.Hubs[chatID]; ok {
		return hub
	}
	n.Hubs[chatID] = models.NewHub()
	n.Hubs[chatID].ChatID=chatID
	return n.Hubs[chatID]
}

func NewArrayRepo() NotificationRepository {
	return &NotificationArrayRepository{Hubs: make(map[uint64]*models.WebSocketHub, 0)}
}



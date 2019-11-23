package usecase

import (
	"github.com/CoolCodeTeam/CoolSupportBackend/notifications/models"
	"github.com/CoolCodeTeam/CoolSupportBackend/notifications/repository"
	chats "github.com/CoolCodeTeam/CoolSupportBackend/chats/usecase"
	users "github.com/CoolCodeTeam/CoolSupportBackend/users/usecase"
)

type NotificationsUseCase interface {
	OpenServerConn(chatID uint64) (*models.WebSocketHub, error)
	OpenClientConn(userID uint64) (*models.WebSocketHub, error)
	SendMessage(chatID uint64, message []byte) error
	HandleCloseConn(chatID uint64) error
}

type NotificationUseCaseImpl struct {
	notificationRepository repository.NotificationRepository
	users users.UsersUseCase
	chats chats.ChatsUseCase
}

func (n *NotificationUseCaseImpl) OpenServerConn(chatID uint64) (*models.WebSocketHub, error) {
	return n.notificationRepository.GetNotificationHub(chatID), nil
}

func (n *NotificationUseCaseImpl) OpenClientConn(userID uint64) (*models.WebSocketHub, error) {
	//get random user_id
	randomID,err:=n.users.GetRandomID()


	//create_chat
	id,err:=n.chats.CreateChat(userID,randomID)
	if err!=nil{
		return &models.WebSocketHub{},err
	}
	//openConn
	n.OpenServerConn(id)
}

func (n *NotificationUseCaseImpl) SendMessage(chatID uint64, message []byte) error {
	hub := n.notificationRepository.GetNotificationHub(chatID)
	if len(hub.Clients) > 0 {
		hub.BroadcastChan <- message
	}
	return nil
}

func (n *NotificationUseCaseImpl) HandleCloseConn(chatID uint64) error{
	//if support still here - email
	hub := n.notificationRepository.GetNotificationHub(chatID)
	if len(hub.Clients) == 0 {
		err:=n.chats.RemoveChat(chatID)
		if err!=nil{
			return err
		}
	}
	return nil
	//if nobody - delete chat
}

func NewNotificationUseCase() NotificationsUseCase {
	return &NotificationUseCaseImpl{notificationRepository: repository.NewArrayRepo()}
}



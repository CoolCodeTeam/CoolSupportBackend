package usecase

import (
	"github.com/CoolCodeTeam/CoolSupportBackend/chats/models"
	"github.com/CoolCodeTeam/CoolSupportBackend/chats/repository"
	users"github.com/CoolCodeTeam/CoolSupportBackend/supports/usecase"
)

type ChatsUseCase interface {
	GetChatsByUserID(ID uint64) ([]models.Chat, error)
	RemoveChat(ID uint64) error
	CreateChat(suppID uint64) (uint64,error)
	GetChat() (uint64 ,error)
}

type ChatsUseCaseImpl struct {
	repository repository.ChatsRepository
	users users.SupportsUseCase
}

func (c *ChatsUseCaseImpl) GetChat() (uint64, error) {
	//get random user_id
	randomID,err:=c.users.GetRandomID()
	if err!=nil{
		return 0,err
	}
	//create_chat
	id,err:=c.CreateChat(randomID)
	return id,err
}

func (c *ChatsUseCaseImpl) RemoveChat(ID uint64) error {
	err:=c.repository.RemoveChat(ID)
	return err
}

func (c *ChatsUseCaseImpl) CreateChat(suppID uint64) (uint64,error){
	id,err:=c.repository.CreateChat(suppID)
	return id,err
}

func (c *ChatsUseCaseImpl) GetChatsByUserID(ID uint64) ([]models.Chat, error) {
	chats, err := c.repository.GetChats(ID)
	var userChats []models.Chat
	if err != nil {
		return chats, err
	}
	return userChats, nil
}

func NewChatsUseCase(repo repository.ChatsRepository) ChatsUseCase {
	return &ChatsUseCaseImpl{
		repository:      repo,
	}
}
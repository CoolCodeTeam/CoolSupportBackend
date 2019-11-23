package repository

import (
	"database/sql"
	"github.com/CoolCodeTeam/CoolSupportBackend/chats/models"
	utils_models "github.com/CoolCodeTeam/CoolSupportBackend/utils/models"
	"net/http"
)

type ChatsRepository interface {
	GetChatByID(ID uint64) (models.Chat, error)
	GetChats(userID uint64) ([]models.Chat, error)
	CreateChat(suppID uint64) (uint64, error)
	RemoveChat(u uint64) error
}

type ChatsDBRepository struct {
	db *sql.DB
}

func (c *ChatsDBRepository) CreateChat(userID uint64) (uint64, error) {
	var chatID uint64
	tx, err := c.db.Begin()
	if err != nil {
		return 0, utils_models.NewServerError(err, http.StatusInternalServerError, "Can not open PutChat transaction "+err.Error())
	}
	defer tx.Rollback()

	_ = c.db.QueryRow("INSERT INTO chats (supportid) VALUES ($1,$2) RETURNING id",
		 userID).Scan(&chatID)
	return chatID, nil

}

func (c *ChatsDBRepository) RemoveChat(u uint64) error {
	_, err := c.db.Exec("DELETE FROM chats WHERE id=$1", u)
	if err != nil {
		return utils_models.NewServerError(err, http.StatusInternalServerError, "Can not delete chat in RemoveChat: "+err.Error())
	}
	return nil
}

func (c ChatsDBRepository) GetChatByID(ID uint64) (models.Chat, error) {
	var result models.Chat

	tx, err := c.db.Begin()
	if err != nil {
		return result, utils_models.NewServerError(err, http.StatusInternalServerError, "can not begin transaction for GetChat: "+err.Error())
	}
	defer tx.Rollback()

	row := tx.QueryRow("SELECT id, supportid FROM chats WHERE id=$1 ", ID)
	err = row.Scan(&result.ID, &result.SupportID)
	if err != nil {
		return result, utils_models.NewClientError(err, http.StatusBadRequest, "chat not exists: "+err.Error())
	}

	return result,nil
}

func (c ChatsDBRepository) GetChats(userID uint64) ([]models.Chat, error) {
	result :=make([]models.Chat,0)

	rows, err := c.db.Query("SELECT id, supportid FROM chats WHERE supportid=$1", userID)
	for rows.Next(){
		var chat models.Chat
		err = rows.Scan(&chat.ID, &chat.SupportID)
		if err != nil {
			return result, utils_models.NewServerError(err, http.StatusBadRequest, "can not GetChats: "+err.Error())
		}
		result = append(result, chat)
	}

	return result,nil
}



func NewChatsDBRepository(db *sql.DB) ChatsRepository {
	return &ChatsDBRepository{db: db}
}


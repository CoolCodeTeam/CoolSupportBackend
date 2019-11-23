package repository

import (
	"database/sql"
	"github.com/CoolCodeTeam/CoolSupportBackend/messages/models"
	utils_models "github.com/CoolCodeTeam/CoolSupportBackend/utils/models"

	"net/http"
	"time"
)

type MessageDBRepository struct {
	DB *sql.DB
}

func (m *MessageDBRepository) PutMessage(message *models.Message) (uint64, error) {
	var chatID uint64
	time, err := time.Parse("02.01.2006 15:04", message.MessageTime)
	if err != nil {
		return 0, utils_models.NewClientError(err, http.StatusBadRequest, "Wrong date format")
	}
	row := m.DB.QueryRow("INSERT into messages (body, chatid,messagetime,isSupport) VALUES ($1,$2,$3,$4,$5,$6) returning id",
		message.Text, message.ChatID, time,message.IsSupp)
	err = row.Scan(&chatID)

	if err != nil {
		return chatID, utils_models.NewServerError(err, http.StatusInternalServerError, "Can not INSERT message in PutMessage "+err.Error())
	}
	return chatID, err

}

func (m *MessageDBRepository) GetMessagesByChatID(chatID uint64) (models.Messages, error) {
	returningMessages := make([]*models.Message, 0)
	rows, err := m.DB.Query("SELECT id,body, chatid,messagetime,isSupport FROM messages where chatid=$1 order by id asc ", chatID)
	if err != nil {
		return models.Messages{}, utils_models.NewServerError(err, http.StatusInternalServerError,
			"Can not get messages in GetMessagesByChatId "+err.Error())
	}
	for rows.Next() {
		var message models.Message
		var messageTime time.Time
		err := rows.Scan(&message.ID, &message.Text, &message.ChatID, &messageTime, &message.IsSupp)

		timeString := messageTime.Format("02.01.2006 15:04")
		message.MessageTime = timeString
		if err != nil {
			return models.Messages{}, utils_models.NewServerError(err, http.StatusInternalServerError,
				"Can not read message in GetMessagesByChatId "+err.Error())
		}
		returningMessages = append(returningMessages, &message)
	}
	return models.Messages{Messages: returningMessages}, nil
}



func NewMessageDbRepository(db *sql.DB) MessageRepository {
	return &MessageDBRepository{DB: db}
}

package delivery

import (
	"encoding/json"
	"github.com/CoolCodeTeam/CoolSupportBackend/messages/models"
	useCase "github.com/CoolCodeTeam/CoolSupportBackend/messages/usecase"
	notifications "github.com/CoolCodeTeam/CoolSupportBackend/notifications/usecase"
	users "github.com/CoolCodeTeam/CoolSupportBackend/supports/usecase"
	utils "github.com/CoolCodeTeam/CoolSupportBackend/utils"
	utils_models "github.com/CoolCodeTeam/CoolSupportBackend/utils/models"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)


type MessageHandlers interface {
	SendMessage(w http.ResponseWriter, r *http.Request)
	GetMessagesByChatID(w http.ResponseWriter, r *http.Request)
}

type MessageHandlersImpl struct {
	Messages      useCase.MessagesUseCase
	Users         users.SupportsUseCase
	Notifications notifications.NotificationsUseCase
	utils         utils.HandlersUtils
}

func (m MessageHandlersImpl) SendMessage(w http.ResponseWriter, r *http.Request) {
	chatID, err := strconv.Atoi(mux.Vars(r)["id"])

	var id uint64
	if err != nil {
		m.utils.LogError(utils_models.NewClientError(err, http.StatusBadRequest, "Bad request: malformed data:("), r)
	}
	if err != nil {
		m.utils.HandleError(err, w, r)
		return
	}
	message, err := parseMessage(r)
	if err != nil {
		m.utils.HandleError(utils_models.NewClientError(err, http.StatusBadRequest, "Bad request: malformed data:("), w, r)
		return
	}
	message.ChatID = uint64(chatID)

		id, err = m.Messages.SaveChatMessage(message)

	if err != nil {
		m.utils.HandleError(err, w, r)
		return
	}
	jsonResponse, err := json.Marshal(map[string]uint64{
		"id": id,
	})
	_, err = w.Write(jsonResponse)
	if err != nil {
		m.utils.LogError(err, r)
	}

	//send to websocket
	message.ID = id
	websocketMessage := models.WebsocketMessage{
		WebsocketEventType: 1,
		Body:               *message,
	}
	websocketJson, err := json.Marshal(websocketMessage)
	if err != nil {
		m.utils.LogError(err, r)
	}
	err = m.Notifications.SendMessage(message.ChatID, websocketJson)
	if err != nil {
		m.utils.LogError(err, r)
	}
}

func (m MessageHandlersImpl) GetMessagesByChatID(w http.ResponseWriter, r *http.Request) {
	var messages models.Messages
	chatID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		m.utils.LogError(utils_models.NewClientError(err, http.StatusBadRequest, "Bad request: malformed data:("), r)
	}
	id, err := m.parseCookie(r)
	if err != nil {
		m.utils.HandleError(err, w, r)
		return
	}

	messages, err = m.Messages.GetChatMessages(uint64(chatID), id)

	if err != nil {
		m.utils.HandleError(err, w, r)
		return
	}
	jsonResponse, err := json.Marshal(messages)
	if err != nil {
		m.utils.HandleError(err, w, r)
	}
	_, err = w.Write(jsonResponse)
	if err != nil {
		m.utils.LogError(err, r)
	}
}

func NewMessageHandlers(useCase useCase.MessagesUseCase, users users.SupportsUseCase,
notificationUseCase notifications.NotificationsUseCase, handlersUtils utils.HandlersUtils) MessageHandlers {
	return &MessageHandlersImpl{
		Messages:      useCase,
		Users:         users,
		Notifications: notificationUseCase,
		utils:         handlersUtils,
	}
}

func (m *MessageHandlersImpl) parseCookie(r *http.Request) (uint64, error) {
	cookie, _ := r.Cookie("session_id")
	id, err := m.Users.GetUserBySession(cookie.Value)
	if err != nil {
		return 0, utils_models.NewClientError(err, http.StatusUnauthorized, "Bad request : not valid cookie:(")
	}
	return id,nil
}

func parseMessage(r *http.Request) (*models.Message, error) {
	var message models.Message
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&message)
	return &message, err
}
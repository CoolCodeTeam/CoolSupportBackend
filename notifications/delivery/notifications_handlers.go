package delivery

import (
	"github.com/CoolCodeTeam/CoolSupportBackend/notifications/usecase"
	users "github.com/CoolCodeTeam/CoolSupportBackend/supports/usecase"
	"github.com/CoolCodeTeam/CoolSupportBackend/utils"
	utils_models "github.com/CoolCodeTeam/CoolSupportBackend/utils/models"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

type NotificationHandlers struct {
	notificationUseCase usecase.NotificationsUseCase
	Users               users.SupportsUseCase
	utils               utils.HandlersUtils
}

func NewNotificationHandlers(notificationUseCase usecase.NotificationsUseCase, users users.SupportsUseCase, utils utils.HandlersUtils) *NotificationHandlers {
	return &NotificationHandlers{notificationUseCase: notificationUseCase, Users: users, utils: utils}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *NotificationHandlers) HandleNewSupportWSConnection(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.utils.HandleError(utils_models.NewServerError(err, http.StatusBadRequest, "Can not upgrade connection"), w, r)
		return
	}

	requestedID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		h.utils.LogError(err, r)
	}

	//Достаем Handler с помощью Messages
	hub, err := h.notificationUseCase.OpenServerConn(uint64(requestedID))
	go hub.Run()
	//Запускаем event loop
	hub.AddClientChan <- ws

	for {
		var m []byte

		_, m, err := ws.ReadMessage()

		if err != nil {
			hub.RemoveClient(ws)
			h.notificationUseCase.HandleCloseConn(uint64(requestedID))
			return
		}
		hub.BroadcastChan <- m
	}

}

func (h *NotificationHandlers) HandleNewClientWSConnection(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.utils.HandleError(utils_models.NewServerError(err, http.StatusBadRequest, "Can not upgrade connection"), w, r)
		return
	}

	//Достаем Handler с помощью Messages
	hub, err := h.notificationUseCase.OpenClientConn()
	go hub.Run()
	//Запускаем event loop
	hub.AddClientChan <- ws

	for {
		var m []byte

		_, m, err := ws.ReadMessage()

		if err != nil {
			hub.RemoveClient(ws)
			err = h.notificationUseCase.HandleCloseConn(hub.ChatID)
			return
		}
		hub.BroadcastChan <- m
	}

}

func (h NotificationHandlers) parseCookie(r *http.Request) (uint64, error) {
	cookie, _ := r.Cookie("session_id")
	ID, err := h.Users.GetUserBySession(cookie.Value)
	if err == nil {
		return ID, nil
	} else {
		return ID, utils_models.NewClientError(nil, http.StatusUnauthorized, "Bad request: no such user :(")
	}
}

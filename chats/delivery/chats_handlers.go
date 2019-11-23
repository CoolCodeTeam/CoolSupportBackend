package delivery

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type ChatHandlers struct {
	Chats    useCase.ChatsUseCase
}

func (c *ChatHandlers) GetChatsByUser(w http.ResponseWriter, r *http.Request) {

	cookie, _ := r.Cookie("session_id")
	cookieID, err := c.Users.GetID(cookie.Value)
	if err != nil {
		c.utils.HandleError(
			models.NewClientError(err, http.StatusUnauthorized, "Bad request : not valid cookie:("),
			w, r)
		return
	}
	if cookieID != uint64(requestedID) {
		c.utils.HandleError(
			models.NewClientError(err, http.StatusUnauthorized, fmt.Sprintf("Actual id: %d, Requested id: %d", cookieID, requestedID)),
			w, r)
		return
	}
	chats, err := c.Chats.GetChatsByUserID(uint64(requestedID))
	if err != nil {
		c.utils.HandleError(err, w, r)
		return
	}
	workspaces, err := c.Chats.GetWorkspacesByUserID(uint64(requestedID))
	if err != nil {
		c.utils.HandleError(err, w, r)
		return
	}
	responseChats := models.ResponseChatsArray{Chats: chats, Workspaces: workspaces}
	jsonChat, err := json.Marshal(responseChats)
	_, err = w.Write(jsonChat)
}


package delivery

import (
	"encoding/json"
	"github.com/CoolCodeTeam/CoolSupportBackend/chats/models"
	"github.com/CoolCodeTeam/CoolSupportBackend/chats/usecase"
	users "github.com/CoolCodeTeam/CoolSupportBackend/supports/usecase"
	"github.com/CoolCodeTeam/CoolSupportBackend/utils"
	utils_models "github.com/CoolCodeTeam/CoolSupportBackend/utils/models"
	"net/http"
)

type ChatHandlers struct {
	Users users.SupportsUseCase
	Chats    usecase.ChatsUseCase
	utils utils.HandlersUtils
}

func (c *ChatHandlers) GetChatsByUser(w http.ResponseWriter, r *http.Request) {

	cookie, _ := r.Cookie("session_id")
	cookieID, err := c.Users.GetUserBySession(cookie.Value)
	if err != nil {
		c.utils.HandleError(
			utils_models.NewClientError(err, http.StatusUnauthorized, "Bad request : not valid cookie:("),
			w, r)
		return
	}

	chats, err := c.Chats.GetChatsByUserID(uint64(cookieID))
	if err != nil {
		c.utils.HandleError(err, w, r)
		return
	}
	responseChats := models.ResponseChatsArray{Chats: chats}
	jsonChat, err := json.Marshal(responseChats)
	_, err = w.Write(jsonChat)
}

func(c *ChatHandlers) GetChat(w http.ResponseWriter, r *http.Request) {
	id,err:=c.Chats.GetChat()
	if err!=nil{
		c.utils.HandleError(err,w,r)
	}
	model:=models.GetChatModel{ID:id}
	response, err := json.Marshal(model)
	_, err = w.Write(response)

}


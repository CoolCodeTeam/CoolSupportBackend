package delivery

import (
	"encoding/json"
	"github.com/CoolCodeTeam/CoolSupportBackend/supports/models"
	"github.com/CoolCodeTeam/CoolSupportBackend/supports/repository"
	"github.com/CoolCodeTeam/CoolSupportBackend/supports/usecase"
	"github.com/CoolCodeTeam/CoolSupportBackend/utils"
	utilsModels "github.com/CoolCodeTeam/CoolSupportBackend/utils/models"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type SupportHandlers struct {
	Supports useCase.SupportsUseCase
	Sessions repository.SessionRepository
	utils    utils.HandlersUtils
}

func (handlers *SupportHandlers) Login(w http.ResponseWriter, r *http.Request) {
	var loginSupport models.Support
	body := r.Body
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&loginSupport)
	if err != nil {
		err = utilsModels.NewClientError(err, http.StatusBadRequest, "Bad request : invalid JSON.")
		handlers.utils.HandleError(err, w, r)
		return
	}

	support, err := handlers.Supports.Login(loginSupport)
	if err != nil {
		handlers.utils.HandleError(err, w, r)
		return
	} else {
		token := uuid.New()
		sessionExpiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{Name: "session_id", Value: token.String(), Expires: sessionExpiration}
		err := handlers.Sessions.Put(cookie.Value, support.ID)
		if err != nil {
			handlers.utils.HandleError(err, w, r)
			return
		}
		support.Password = ""
		body, err := json.Marshal(support)
		if err != nil {
			handlers.utils.HandleError(err, w, r)
			return
		}
		http.SetCookie(w, &cookie)
		w.Header().Set("content-type", "application/json")

		_, err = w.Write(body)
		if err != nil {
			handlers.utils.HandleError(err, w, r)
			return
		}
		return
	}
}

func (handlers *SupportHandlers) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Cookie("session_id")
	err := handlers.Sessions.Remove(session.Value)
	if err != nil {
		handlers.utils.HandleError(
			utilsModels.NewClientError(err, http.StatusUnauthorized, "Bad request : not valid cookie:("),
			w, r)
	}
	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)
}

func (handlers *SupportHandlers) GetSupportBySession(w http.ResponseWriter, r *http.Request) {
	sessionID, err := r.Cookie("session_id")
	if err != nil {
		handlers.utils.HandleError(utilsModels.NewClientError(err, http.StatusUnauthorized, "Not authorized:("), w, r)
		return
	}

	support, err := handlers.parseCookie(sessionID)
	if err != nil {
		handlers.utils.HandleError(err, w, r)
		return
	}

	body, err := json.Marshal(support)
	if err != nil {
		handlers.utils.HandleError(err, w, r)
		return
	}

	_, err = w.Write(body)
	if err != nil {
		handlers.utils.HandleError(err, w, r)
		return
	}

}

func (handlers *SupportHandlers) parseCookie(cookie *http.Cookie) (models.Support, error) {
	id, err := handlers.Sessions.GetID(cookie.Value)
	if err != nil {
		return models.Support{}, utilsModels.NewClientError(err, http.StatusUnauthorized, "Bad request : not valid cookie:(")
	}
	support, err := handlers.Supports.GetSupportByID(id)
	if err == nil {
		return support, nil
	} else {
		return support, err
	}
}

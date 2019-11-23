package useCase

import (
	"github.com/CoolCodeTeam/CoolSupportBackend/supports/models"
	"github.com/CoolCodeTeam/CoolSupportBackend/supports/repository"
	utilsModels "github.com/CoolCodeTeam/CoolSupportBackend/utils/models"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

type SupportsUseCase interface {
	GetSupportByID(id uint64) (models.Support, error)
	Login(support models.Support) (models.Support, error)
	GetUserBySession(session string) (uint64, error)
	GetRandomID() (uint64, error)
}

type supportUseCase struct {
	repository repository.SupportRepo
	sessions   repository.SessionRepository
}

func (u *supportUseCase) Login(loginSupport models.Support) (models.Support, error) {
	support, err := u.repository.GetSupportByEmail(loginSupport.Email)
	if err != nil {
		err = utilsModels.NewClientError(nil, http.StatusBadRequest, "Bad request: malformed data")
		return support, err
	}

	if comparePasswords(support.Password, loginSupport.Password) {
		return support, nil
	} else {
		err = utilsModels.NewClientError(nil, http.StatusBadRequest, "Bad request: wrong password")
		return support, err
	}

}

func NewSupportUseCase(repo repository.SupportRepo, sessions repository.SessionRepository) SupportsUseCase {
	return &supportUseCase{
		repository: repo,
		sessions:   sessions,
	}
}

func (u *supportUseCase) GetSupportByID(id uint64) (models.Support, error) {
	support, err := u.repository.GetSupportByID(id)
	if err != nil {
		return support, err
	}
	if !u.Valid(support) {
		return support, utilsModels.NewClientError(nil, http.StatusUnauthorized, "Bad request: no such support :(")
	}
	return support, nil
}

func (u *supportUseCase) GetRandomID() (uint64, error) {
	randID, err := u.repository.GetRandomID()
	if err != nil {
		return 0, err
	}
	return randID, nil
}

func (u *supportUseCase) Valid(support models.Support) bool {
	return support.Email != ""
}

func comparePasswords(hashedPassword string, plainPassword string) bool {
	return hashedPassword == plainPassword

}

func (u *supportUseCase) GetUserBySession(session string) (uint64, error) {
	id, err := u.sessions.GetID(session)
	return id, err
}

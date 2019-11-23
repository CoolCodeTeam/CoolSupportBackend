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
}

type supportUseCase struct {
	repository repository.SupportRepo
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

func NewSupportUseCase(repo repository.SupportRepo) SupportsUseCase {
	return &supportUseCase{
		repository: repo,
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

func (u *supportUseCase) Valid(support models.Support) bool {
	return support.Email != ""
}

func comparePasswords(hashedPassword string, plainPassword string) bool {
	byteHash := []byte(hashedPassword)
	err := bcrypt.CompareHashAndPassword(byteHash, []byte(plainPassword))
	if err != nil {
		return false
	}
	return true
}

//TODO: use on support creation
func hashAndSalt(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}
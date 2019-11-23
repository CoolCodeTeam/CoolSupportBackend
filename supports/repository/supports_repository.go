package repository

import "github.com/CoolCodeTeam/CoolSupportBackend/supports/models"

type SupportRepo interface {
	GetSupportByEmail(email string) (models.Support, error)
	GetSupportByID(ID uint64) (models.Support, error)
}

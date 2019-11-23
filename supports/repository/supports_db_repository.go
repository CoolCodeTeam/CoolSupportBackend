package repository

import (
	"database/sql"
	"github.com/CoolCodeTeam/CoolSupportBackend/supports/models"
	utilsModels "github.com/CoolCodeTeam/CoolSupportBackend/utils/models"
	"net/http"
)

type DBSupportStore struct {
	DB *sql.DB
}

func (SupportStore *DBSupportStore) GetSupportByID(ID uint64) (models.Support, error) {
	support := &models.Support{}
	selectStr := "SELECT id, email, password FROM supports WHERE id = $1"
	row := SupportStore.DB.QueryRow(selectStr, ID)

	err := row.Scan(&support.ID, &support.Email, &support.Password)
	if err != nil {
		return *support, utilsModels.NewServerError(err, http.StatusInternalServerError, "Can not get support: "+err.Error())
	}
	return *support, nil
}

func (SupportStore *DBSupportStore) GetSupportByEmail(email string) (models.Support, error) {
	support := &models.Support{}
	selectStr := "SELECT id, email, password FROM supports WHERE email = $1"
	row := SupportStore.DB.QueryRow(selectStr, email)

	err := row.Scan(&support.ID, &support.Email, &support.Password)

	if err != nil {
		return *support, utilsModels.NewServerError(err, http.StatusInternalServerError, "Can not get support: "+err.Error())
	}
	return *support, nil
}

func NewSupportDBStore(db *sql.DB) SupportRepo {
	return &DBSupportStore{
		db,
	}
}

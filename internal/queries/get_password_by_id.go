package queries

import (
	"errors"
	"fmt"

	"github.com/iskorotkov/passwordmanager/internal/database"
	"github.com/iskorotkov/passwordmanager/internal/models"
	"gorm.io/gorm"
)

var (
	ErrGetPasswordByIDInternalError = fmt.Errorf("internal error")
	ErrGetPasswordByIDNotFound      = fmt.Errorf("password not found")
)

func GetPasswordByID(id uint, p *models.Password) database.Query {
	return func(db *gorm.DB) error {
		err := db.First(p, id).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrGetPasswordByIDNotFound
		}

		if err != nil {
			return ErrGetPasswordByIDInternalError
		}

		return nil
	}
}

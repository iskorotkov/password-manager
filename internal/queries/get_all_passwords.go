package queries

import (
	"fmt"

	"github.com/iskorotkov/passwordmanager/internal/database"
	"github.com/iskorotkov/passwordmanager/internal/models"
	"gorm.io/gorm"
)

var ErrGetAllPasswordsInternalError = fmt.Errorf("internal error")

func GetAllPasswords(p *[]models.Password) database.Query {
	return func(db *gorm.DB) error {
		err := db.Find(p).Error
		if err != nil {
			return ErrGetAllPasswordsInternalError
		}

		return nil
	}
}

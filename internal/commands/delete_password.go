package commands

import (
	"errors"
	"fmt"

	"github.com/iskorotkov/password-manager/internal/database"
	"github.com/iskorotkov/password-manager/internal/models"
	"gorm.io/gorm"
)

var (
	ErrDeletePasswordInternalError = fmt.Errorf("internal error")
	ErrDeletePasswordNotFound      = fmt.Errorf("password not found")
)

func DeletePassword(id uint) database.Command {
	return func(db *gorm.DB) error {
		tx := db.Delete(&models.Password{}, id) //nolint:exhaustivestruct

		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return ErrDeletePasswordNotFound
		}

		if tx.Error != nil {
			return ErrDeletePasswordInternalError
		}

		if tx.RowsAffected == 0 {
			return ErrDeletePasswordNotFound
		}

		return nil
	}
}

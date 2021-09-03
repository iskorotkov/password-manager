package commands

import (
	"errors"
	"fmt"

	"github.com/iskorotkov/passwordmanager/internal/database"
	"github.com/iskorotkov/passwordmanager/internal/models"
	"gorm.io/gorm"
)

var (
	ErrDeletePasswordInternalError = fmt.Errorf("internal error")
	ErrDeletePasswordNotFound      = fmt.Errorf("password not found")
)

func DeletePassword(id uint) database.Command {
	return func(db *gorm.DB) error {
		err := db.Delete(&models.Password{}, id).Error //nolint:exhaustivestruct
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrDeletePasswordNotFound
		}

		if err != nil {
			return ErrDeletePasswordInternalError
		}

		return nil
	}
}

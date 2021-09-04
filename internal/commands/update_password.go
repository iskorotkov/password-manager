package commands

import (
	"errors"
	"fmt"

	"github.com/iskorotkov/password-manager/internal/database"
	"github.com/iskorotkov/password-manager/internal/models"
	"gorm.io/gorm"
)

var (
	ErrUpdatePasswordInternalError = fmt.Errorf("internal error")
	ErrUpdatePasswordNotFound      = fmt.Errorf("not found")
)

func UpdatePassword(p models.Password) database.Command {
	return func(db *gorm.DB) error {
		err := db.Updates(&p).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUpdatePasswordNotFound
		}

		if err != nil {
			return ErrUpdatePasswordInternalError
		}

		return nil
	}
}

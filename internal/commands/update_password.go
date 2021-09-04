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
		// Force deleted_at check as otherwise it would be possible to update deleted password
		// and thus get RowsAffected > 0 and HTTP status code of 200, which isn't correct.
		tx := db.Where("deleted_at is null").Updates(&p)

		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return ErrUpdatePasswordNotFound
		}

		if tx.Error != nil {
			return ErrUpdatePasswordInternalError
		}

		if tx.RowsAffected == 0 {
			return ErrUpdatePasswordNotFound
		}

		return nil
	}
}

package commands

import (
	"fmt"

	"github.com/iskorotkov/passwordmanager/internal/database"
	"github.com/iskorotkov/passwordmanager/internal/models"
	"gorm.io/gorm"
)

var (
	ErrUpdatePasswordInternalError = fmt.Errorf("internal error")
	ErrUpdatePasswordNotFound      = fmt.Errorf("not found")
)

func UpdatePassword(p models.Password) database.Command {
	return func(db *gorm.DB) error {
		err := db.Updates(&p).Error
		if err == gorm.ErrRecordNotFound {
			return ErrUpdatePasswordNotFound
		}

		if err != nil {
			return ErrUpdatePasswordInternalError
		}

		return nil
	}
}

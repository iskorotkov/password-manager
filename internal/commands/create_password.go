package commands

import (
	"fmt"

	"github.com/iskorotkov/password-manager/internal/database"
	"github.com/iskorotkov/password-manager/internal/models"
	"gorm.io/gorm"
)

var ErrCreatePasswordInternalError = fmt.Errorf("internal error")

func CreatePassword(p *models.Password) database.Command {
	return func(db *gorm.DB) error {
		err := db.Create(p).Error
		if err != nil {
			return ErrCreatePasswordInternalError
		}

		return nil
	}
}

package postgres

import (
	"fmt"

	"github.com/iskorotkov/password-manager/internal/database"
	"github.com/iskorotkov/password-manager/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormDB struct {
	db *gorm.DB
}

func New(connString string) (GormDB, error) {
	db, err := gorm.Open(postgres.Open(connString))
	if err != nil {
		return GormDB{}, fmt.Errorf("error connecting to database: %w", err)
	}

	err = db.AutoMigrate(
		&models.Password{}, //nolint:exhaustivestruct
	)
	if err != nil {
		return GormDB{}, fmt.Errorf("error migrating models: %w", err)
	}

	return GormDB{db: db}, nil
}

func (d GormDB) Query(q database.Query) error {
	return q(d.db)
}

func (d GormDB) Exec(c database.Command) error {
	return c(d.db)
}

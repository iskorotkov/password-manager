package postgres

import (
	"fmt"
	"log"
	"time"

	"github.com/iskorotkov/password-manager/internal/database"
	"github.com/iskorotkov/password-manager/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	retries  = 5
	interval = time.Second * 6
)

type GormDB struct {
	db *gorm.DB
}

func New(connString string) (GormDB, error) {
	var (
		db  *gorm.DB
		err error
	)

	for i := 0; i < retries; i++ {
		db, err = gorm.Open(postgres.Open(connString))
		if err == nil {
			break
		}

		log.Printf("error connecting to database, will wait for %v and repeat (%d/%d)", interval, i+1, retries)

		time.Sleep(interval)
	}

	if err != nil {
		return GormDB{}, fmt.Errorf("error connecting to database: %w", err)
	}

	log.Printf("connected to DB")

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

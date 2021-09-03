package database

import (
	"gorm.io/gorm"
)

type DB interface {
	Query(q Query) error
	Exec(c Command) error
}

type Query func(db *gorm.DB) error

type Command func(db *gorm.DB) error

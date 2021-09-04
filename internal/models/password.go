package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	openapi "github.com/iskorotkov/password-manager/go"
	"gorm.io/gorm"
)

type Password struct {
	gorm.Model
	Website  string
	Username string
	Password string
}

func (p Password) Validate() error {
	return validation.ValidateStruct(&p, //nolint:wrapcheck
		validation.Field(&p.Website, validation.Required, is.URL),
		validation.Field(&p.Username, validation.Required, is.Alphanumeric),
		validation.Field(&p.Password, validation.Required),
	)
}

func (p Password) ToDTO() openapi.Password {
	return openapi.Password{
		Id:       int32(p.ID),
		Website:  p.Website,
		Username: p.Username,
		Password: p.Password,
	}
}

func (p Password) FromDTO(password openapi.Password) Password {
	return Password{
		Model:    gorm.Model{ID: uint(password.Id)}, //nolint:exhaustivestruct
		Website:  password.Website,
		Username: password.Username,
		Password: password.Password,
	}
}

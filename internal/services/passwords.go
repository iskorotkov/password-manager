package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	openapi "github.com/iskorotkov/password-manager/go"
	"github.com/iskorotkov/password-manager/internal/commands"
	"github.com/iskorotkov/password-manager/internal/database"
	"github.com/iskorotkov/password-manager/internal/models"
	"github.com/iskorotkov/password-manager/internal/queries"
)

var (
	ErrIDsDoNotMatch = fmt.Errorf("id in path and body don't match")
	ErrNonPositiveID = fmt.Errorf("invalid password id")
)

type PasswordService struct {
	db database.DB
}

func (p PasswordService) ApiV1PasswordsGet(_ context.Context) (openapi.ImplResponse, error) { //nolint:revive,stylecheck
	var passwords []models.Password

	err := p.db.Query(queries.GetAllPasswords(&passwords))
	if err != nil {
		return openapi.ImplResponse{
			Code: http.StatusInternalServerError,
			Body: nil,
		}, err //nolint:wrapcheck
	}

	results := make([]openapi.Password, 0, len(passwords))
	for _, p := range passwords {
		results = append(results, p.ToDTO())
	}

	return openapi.ImplResponse{
		Code: http.StatusOK,
		Body: results,
	}, nil
}

func (p PasswordService) ApiV1PasswordsIdDelete( //nolint:revive,stylecheck
	_ context.Context,
	i int32,
) (openapi.ImplResponse, error) {
	if i <= 0 {
		return openapi.ImplResponse{
			Code: http.StatusBadRequest,
			Body: nil,
		}, ErrNonPositiveID
	}

	err := p.db.Exec(commands.DeletePassword(uint(i)))
	if errors.Is(err, commands.ErrDeletePasswordNotFound) {
		return openapi.ImplResponse{
			Code: http.StatusNotFound,
			Body: nil,
		}, err //nolint:wrapcheck
	}

	if err != nil {
		return openapi.ImplResponse{
			Code: http.StatusInternalServerError,
			Body: nil,
		}, err //nolint:wrapcheck
	}

	return openapi.ImplResponse{
		Code: http.StatusOK,
		Body: "password deleted",
	}, nil
}

func (p PasswordService) ApiV1PasswordsIdGet( //nolint:revive,stylecheck
	_ context.Context,
	i int32,
) (openapi.ImplResponse, error) {
	if i <= 0 {
		return openapi.ImplResponse{
			Code: http.StatusBadRequest,
			Body: nil,
		}, ErrNonPositiveID
	}

	var password models.Password

	err := p.db.Query(queries.GetPasswordByID(uint(i), &password))
	if errors.Is(err, queries.ErrGetPasswordByIDNotFound) {
		return openapi.ImplResponse{
			Code: http.StatusNotFound,
			Body: nil,
		}, err //nolint:wrapcheck
	}

	if err != nil {
		return openapi.ImplResponse{
			Code: http.StatusInternalServerError,
			Body: nil,
		}, err //nolint:wrapcheck
	}

	return openapi.ImplResponse{
		Code: http.StatusOK,
		Body: password.ToDTO(),
	}, nil
}

func (p PasswordService) ApiV1PasswordsIdPut( //nolint:revive,stylecheck
	_ context.Context,
	i int32,
	password openapi.Password,
) (openapi.ImplResponse, error) {
	if password.Id != 0 && i != password.Id {
		return openapi.ImplResponse{
			Code: http.StatusBadRequest,
			Body: nil,
		}, ErrIDsDoNotMatch
	}

	if i <= 0 {
		return openapi.ImplResponse{
			Code: http.StatusBadRequest,
			Body: nil,
		}, ErrNonPositiveID
	}

	password.Id = i

	model := models.Password{}.FromDTO(password) //nolint:exhaustivestruct

	err := model.Validate()
	if err != nil {
		return openapi.ImplResponse{
			Code: http.StatusBadRequest,
			Body: nil,
		}, err //nolint:wrapcheck
	}

	err = p.db.Exec(commands.UpdatePassword(model))
	if errors.Is(err, commands.ErrUpdatePasswordNotFound) {
		return openapi.ImplResponse{
			Code: http.StatusNotFound,
			Body: nil,
		}, err //nolint:wrapcheck
	}

	if err != nil {
		return openapi.ImplResponse{
			Code: http.StatusInternalServerError,
			Body: nil,
		}, err //nolint:wrapcheck
	}

	return openapi.ImplResponse{
		Code: http.StatusOK,
		Body: "password updated",
	}, nil
}

func (p PasswordService) ApiV1PasswordsPost( //nolint:revive,stylecheck
	_ context.Context,
	password openapi.Password,
) (openapi.ImplResponse, error) {
	password.Id = 0

	model := models.Password{}.FromDTO(password) //nolint:exhaustivestruct

	err := model.Validate()
	if err != nil {
		return openapi.ImplResponse{
			Code: http.StatusBadRequest,
			Body: nil,
		}, err //nolint:wrapcheck
	}

	err = p.db.Exec(commands.CreatePassword(&model))
	if err != nil {
		return openapi.ImplResponse{
			Code: http.StatusInternalServerError,
			Body: nil,
		}, err //nolint:wrapcheck
	}

	return openapi.ImplResponse{
		Code: http.StatusCreated,
		Body: model.ToDTO(),
	}, nil
}

func NewPasswordService(db database.DB) PasswordService {
	return PasswordService{db: db}
}

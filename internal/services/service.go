package services

import (
	"context"
	"fmt"
	"net/http"

	openapi "github.com/iskorotkov/passwordmanager/go"
	"github.com/iskorotkov/passwordmanager/internal/commands"
	"github.com/iskorotkov/passwordmanager/internal/database"
	"github.com/iskorotkov/passwordmanager/internal/models"
	"github.com/iskorotkov/passwordmanager/internal/queries"
)

type PasswordService struct {
	db database.DB
}

func (p PasswordService) ApiV1PasswordsGet(_ context.Context) (openapi.ImplResponse, error) {
	var passwords []models.Password

	err := p.db.Query(queries.GetAllPasswords(&passwords))
	if err != nil {
		return openapi.ImplResponse{
			Code: http.StatusInternalServerError,
			Body: err,
		}, err
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

func (p PasswordService) ApiV1PasswordsIdDelete(_ context.Context, i int32) (openapi.ImplResponse, error) {
	err := p.db.Exec(commands.DeletePassword(uint(i)))
	if err == commands.ErrDeletePasswordNotFound {
		return openapi.ImplResponse{
			Code: http.StatusNotFound,
			Body: err,
		}, err
	}

	if err != nil {
		return openapi.ImplResponse{
			Code: http.StatusInternalServerError,
			Body: err,
		}, err
	}

	return openapi.ImplResponse{
		Code: http.StatusOK,
		Body: "password deleted",
	}, nil
}

func (p PasswordService) ApiV1PasswordsIdGet(_ context.Context, i int32) (openapi.ImplResponse, error) {
	var password models.Password

	err := p.db.Query(queries.GetPasswordByID(uint(i), &password))
	if err == queries.ErrGetPasswordByIDNotFound {
		return openapi.ImplResponse{
			Code: http.StatusNotFound,
			Body: err,
		}, err
	}

	if err != nil {
		return openapi.ImplResponse{
			Code: http.StatusInternalServerError,
			Body: err,
		}, err
	}

	return openapi.ImplResponse{
		Code: http.StatusOK,
		Body: password.ToDTO(),
	}, nil
}

func (p PasswordService) ApiV1PasswordsIdPut(_ context.Context, i int32, password openapi.Password) (openapi.ImplResponse, error) {
	if password.Id != 0 && i != password.Id {
		err := fmt.Errorf("id in path and body doesn't match")

		return openapi.ImplResponse{
			Code: http.StatusBadRequest,
			Body: err,
		}, err
	}

	password.Id = i

	model := models.Password{}.FromDTO(password)

	err := model.Validate()
	if err != nil {
		return openapi.ImplResponse{
			Code: http.StatusBadRequest,
			Body: err,
		}, err
	}

	err = p.db.Exec(commands.UpdatePassword(model))
	if err == commands.ErrUpdatePasswordNotFound {
		return openapi.ImplResponse{
			Code: http.StatusNotFound,
			Body: err,
		}, err
	}

	if err != nil {
		return openapi.ImplResponse{
			Code: http.StatusInternalServerError,
			Body: err,
		}, err
	}

	return openapi.ImplResponse{
		Code: http.StatusOK,
		Body: "password updated",
	}, nil
}

func (p PasswordService) ApiV1PasswordsPost(_ context.Context, password openapi.Password) (openapi.ImplResponse, error) {
	password.Id = 0

	model := models.Password{}.FromDTO(password)

	err := model.Validate()
	if err != nil {
		return openapi.ImplResponse{
			Code: http.StatusBadRequest,
			Body: err,
		}, err
	}

	err = p.db.Exec(commands.CreatePassword(&model))
	if err != nil {
		return openapi.ImplResponse{
			Code: http.StatusInternalServerError,
			Body: err,
		}, err
	}

	return openapi.ImplResponse{
		Code: http.StatusCreated,
		Body: model.ToDTO(),
	}, nil
}

func NewPasswordService(db database.DB) PasswordService {
	return PasswordService{db: db}
}

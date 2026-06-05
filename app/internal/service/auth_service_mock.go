package service

import (
	"context"
	"errors"

	"github.com/Chocobo11218/go-auth-jwt/app/internal/model"
	//"github.com/stretchr/testify/mock"
)

type authServiceMock struct {
	//mock.Mock
}

func NewAuthServiceMock() *authServiceMock {
	return &authServiceMock{}
}

// localhost:4001/api/v1/register
func (m *authServiceMock) Register(_ context.Context, req *model.RegisterRequest) (model.AppResponse, error) {

	switch req.Email {
	// email already exist
	case "taken@example.com":
		return model.AppResponse{}, errors.New(model.EmailAlreadyExistMessage)
	// service unavailable
	case "off@example.com":
		return model.AppResponse{}, errors.New(model.ServiceUnavailableMessage)
	}
	// register success
	return model.AppResponse{
		Code:    model.StatusSuccess,
		Message: "Register success",
	}, nil

}

func (m *authServiceMock) Login(_ context.Context, req *model.LoginRequest) (model.AppResponse, error) {

	switch req.Email {
	// wrong password or email
	case "wrong@example.com":
		return model.AppResponse{}, errors.New(model.InvalidCredentialMessage)

	case "off@example.com":
		return model.AppResponse{}, errors.New(model.ServiceUnavailableMessage)
	}
	// success
	return model.AppResponse{
		Code:    model.StatusSuccess,
		Message: "Success",
		Data: model.TokenData{
			AccessToken: "mock-jwt-token",
		},
	}, nil
}

/*
{
  "email": "customer@example.com",
  "password": "12345678",
  "first_name": "John",
  "last_name": "Doe",
  "phone_number": "0812345678"
}

{
  "email": "customer@example.com",
  "password": "12345678"
}
*/

package service

import (
	"context"
	"errors"

	"github.com/Chocobo11218/go-auth-jwt/app/internal/model"
	"github.com/stretchr/testify/mock"
)

type authServiceMock struct {
	mock.Mock
}

func NewAuthServiceMock() *authServiceMock {
	return &authServiceMock{}
}

func (m *authServiceMock) Register(ctx context.Context, req *model.RegisterRequest) (model.AppResponse, error) {
	// email already exist
	if req.Email == "taken@example.com" {
		return model.AppResponse{}, errors.New("Email already exist")
	}
	return model.AppResponse{}, nil
}

func (m *authServiceMock) Login(ctx context.Context, req *model.LoginRequest) (model.TokenData, error) {
	// wrong password
	if req.Password != "123456" {
		return model.TokenData{}, errors.New("Invalid email or password")
	}
	return model.TokenData{AccessToken: "mock-jwt-token"}, nil
}

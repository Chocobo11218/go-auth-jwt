package service

import (
	"context"
	"errors"

	"github.com/Chocobo11218/go-auth-jwt/app/internal/model"
	"github.com/Chocobo11218/go-auth-jwt/app/internal/repository"
	"github.com/Chocobo11218/go-auth-jwt/app/pkg/logger"
	"go.uber.org/zap"
)

type UserService interface {
	GetProfile(ctx context.Context, userID uint) (*model.AppResponse, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetProfile(ctx context.Context, userID uint) (*model.AppResponse, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		logger.Error(ctx, "UserService - GetProfile: failed to get user", zap.Error(err))
		return nil, err
	}
	if user == nil {
		logger.Info(ctx, "UserService - GetProfile: user not found", zap.Uint("user_id", userID))
		return nil, errors.New(model.UserNotFoundMessage)
	}

	return &model.AppResponse{
		Code:    model.StatusSuccess,
		Message: model.GetMeSuccessMessage,
		Data: model.ProfileResponse{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		},
	}, nil
}

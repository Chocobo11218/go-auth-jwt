package service

import (
	"context"
	"errors"
	"time"

	"github.com/Chocobo11218/go-auth-jwt/app/internal/model"
	"github.com/Chocobo11218/go-auth-jwt/app/internal/repository"
	"github.com/Chocobo11218/go-auth-jwt/app/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// AuthService is the port (interface) that the HTTP handler depends on.
// Any adapter (HTTP, gRPC, CLI) talks to this, never directly to the repository.
type AuthService interface {
	Register(ctx context.Context, req *model.RegisterRequest) (model.AppResponse, error)
	Login(ctx context.Context, req *model.LoginRequest) (model.AppResponse, error)
}

type authService struct {
	userRepo     repository.UserRepository
	jwtSecretKey string
	jwtTTL       time.Duration
}

// constructor
func NewAuthService(
	userRepo repository.UserRepository,
	jwtSecretKey string,
	jwtTTL time.Duration,
) AuthService {
	return &authService{
		userRepo:     userRepo,
		jwtSecretKey: jwtSecretKey,
		jwtTTL:       jwtTTL,
	}
}

// creates a new user after validating that the email is not already taken
func (s *authService) Register(ctx context.Context, req *model.RegisterRequest) (model.AppResponse, error) {

	// check service hours
	if err := s.checkServiceHours(); err != nil {
		return model.AppResponse{}, errors.New(model.ServiceUnavailableMessage)
	}

	// check for duplicate email
	exists, err := s.userRepo.ExistByEmail(ctx, req.Email)

	logger.Info(ctx, "AuthService - Register ExistByEmail called")

	if err != nil {
		logger.Error(ctx, "AuthService - Register failed to check email existence",
            zap.String("email", req.Email),
            zap.Error(err),
        )
		return model.AppResponse{}, errors.New(model.GenericErrorMessage)
	}
	if exists {
		logger.Error(ctx, "AuthService - Register email already exists",
            zap.String("email", req.Email),
        )
		return model.AppResponse{}, errors.New(model.EmailAlreadyExistMessage)
	}

	// hash the password before storing
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error(ctx, "AuthService - Register failed to hash password", zap.Error(err))
		return model.AppResponse{}, errors.New(model.GenericErrorMessage)
	}

	saveUserDB := &model.User{
		Email:       req.Email,
		Password:    string(hashPassword),
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
	}

	err = s.userRepo.CreateUser(ctx, saveUserDB)
	if err != nil {
		logger.Error(ctx, "AuthService - Register failed to create user",
            zap.String("email", req.Email),
            zap.Error(err),
        )
		return model.AppResponse{}, errors.New(model.GenericErrorMessage)
	}

	registerResponse := model.AppResponse{
		Code:    model.StatusSuccess,
		Message: "Register Success",
	}

	return registerResponse, nil
}

// gives a registered user an access token if the credentials are valid
func (s *authService) Login(ctx context.Context, req *model.LoginRequest) (model.AppResponse, error) {

	// check service hours
	if err := s.checkServiceHours(); err != nil {
		return model.AppResponse{}, errors.New(model.ServiceUnavailableMessage)
	}

	// find user email
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	
	logger.Info(ctx, "AuthService - Login GetUserByEmail called")
	
	if err != nil {
		logger.Error(ctx, "AuthService - Login failed to find user",
            zap.Error(err),
        )
		return model.AppResponse{}, errors.New(model.GenericErrorMessage)
	}
	if user == nil {
		logger.Error(ctx, "AuthService - Login user not found",)
		return model.AppResponse{}, errors.New(model.InvalidCredentialMessage)
	}

	// compare user input password against stored hashed password
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	)
	if err != nil {
		logger.Error(ctx, "AuthService - Login invalid password", zap.Error(err))
		return model.AppResponse{}, errors.New(model.InvalidCredentialMessage)
	}

	// pass = return jwt
	// create jwt token
	token, err := s.generateJWT(user.ID)
	if err != nil {
		logger.Error(ctx, "AuthService - Login failed to generate jwt", zap.Error(err))
		return model.AppResponse{}, errors.New(model.GenericErrorMessage)
	}

	response := model.AppResponse{
		Code:    model.StatusSuccess,
		Message: "Login Success",
		Data:    model.TokenData{AccessToken: token},
	}
	return response, nil
}

func (s *authService) checkServiceHours() error {
	hour := time.Now().Hour()
	if hour < 6 || hour >= 23 {
		return errors.New(model.ServiceUnavailableMessage)
	}
	return nil
}

func (s *authService) generateJWT(userID uint) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(s.jwtTTL).Unix()

	t, err := token.SignedString([]byte(s.jwtSecretKey))
	if err != nil {
		return "", err
	}
	return t, nil
}

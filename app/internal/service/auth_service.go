package service

import (
	"context"
	"errors"
	"time"

	"github.com/Chocobo11218/go-auth-jwt/app/internal/config"
	"github.com/Chocobo11218/go-auth-jwt/app/internal/model"
	"github.com/Chocobo11218/go-auth-jwt/app/internal/repository"
	"github.com/Chocobo11218/go-auth-jwt/app/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, req *model.RegisterRequest) (*model.AppResponse, error)
	Login(ctx context.Context, req *model.LoginRequest) (*model.AppResponse, error)
}

type authService struct {
	userRepo repository.UserRepository
	conf     *config.AppConfig
}

// constructor
func NewAuthService(
	userRepo repository.UserRepository,
	conf *config.AppConfig,
) AuthService {
	return &authService{
		userRepo: userRepo,
		conf:     conf,
	}
}

// creates a new user after validating that the email is not already taken
func (s *authService) Register(ctx context.Context, req *model.RegisterRequest) (*model.AppResponse, error) {

	// check service hours
	if err := checkServiceHours(); err != nil {
		logger.Info(ctx, "Register Service: Service hours restricted",
			zap.Error(err),
		)
		return nil, errors.New(model.ServiceUnavailableMessage)
	}

	// check for duplicate email
	exists, err := s.userRepo.ExistByEmail(ctx, req.Email)

	if err != nil {
		logger.Info(ctx, "Register Service: Failed to check email existence",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return nil, err
	}
	if exists {
		logger.Info(ctx, "Register Service: Email already exists",
			zap.String("email", req.Email),
		)
		return nil, errors.New(model.EmailAlreadyExistMessage)
	}

	// hash the password before storing
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error(ctx, "Register Service: Failed to hash password",
			zap.Error(err),
		)
		return nil, err
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
		logger.Error(ctx, "Register Service: Failed to create user",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return nil, err
	}

	registerResponse := &model.AppResponse{
		Code:    model.StatusSuccess,
		Message: model.RegisterSuccessMessage,
	}

	return registerResponse, nil
}

// gives a registered user an access token if the credentials are valid
func (s *authService) Login(ctx context.Context, req *model.LoginRequest) (*model.AppResponse, error) {

	// check service hours
	if err := checkServiceHours(); err != nil {
		logger.Info(ctx, "Login Service: Service hours restricted",
			zap.Error(err),
		)
		return nil, errors.New(model.ServiceUnavailableMessage)
	}

	// find user email
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		logger.Error(ctx, "Login Service: User look up failed",
			zap.Error(err),
		)
		return nil, err
	}
	// no user found
	if user == nil {
		logger.Error(ctx, "Login Service: User not found")
		return nil, errors.New(model.InvalidCredentialMessage)
	}

	// compare user input password against stored hashed password
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	)
	if err != nil {
		logger.Error(ctx, "Login Service: Invalid password",
			zap.Error(err),
		)
		return nil, errors.New(model.InvalidCredentialMessage)
	}

	// pass = create and return jwt token
	token, err := s.generateJWT(user.ID)
	if err != nil {
		logger.Error(ctx, "Login Service: jwt generation failed",
			zap.Error(err),
		)
		return nil, err
	}

	response := &model.AppResponse{
		Code:    model.StatusSuccess,
		Message: model.LoginSuccessMessage,
		Data:    model.TokenData{
			AccessToken: token,
		},
	}
	return response, nil
}

func checkServiceHours() error {
	hour := time.Now().In(time.Local).Hour()
	if hour < 6 || hour >= 23 {
		return errors.New(model.ServiceUnavailableMessage)
	}
	return nil
}

func (s *authService) generateJWT(userID uint) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["jti"] = uuid.New().String()
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(s.conf.JWT.Expire).Unix()

	t, err := token.SignedString([]byte(s.conf.Secret.JWTSecret))
	if err != nil {
		return "", err
	}
	return t, nil
}

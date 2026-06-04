package service

import (
	"context"

	//"encoding/json"
	"errors"
	//"fmt"
	"strconv"
	"time"

	"github.com/Chocobo11218/go-auth-jwt/app/internal/model"
	"github.com/Chocobo11218/go-auth-jwt/app/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthService is the port (interface) that the HTTP handler depends on.
// Any adapter (HTTP, gRPC, CLI) talks to this, never directly to the repository.
type AuthService interface {
	Register(ctx context.Context, req *model.RegisterRequest) (*model.AppResponse, error)
	Login(ctx context.Context, req *model.LoginRequest) (*model.AppResponse, error)
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
func (s *authService) Register(ctx context.Context, req *model.RegisterRequest) (*model.AppResponse, error) {

	// check service hours
	if err := s.checkServiceHours(); err != nil {
		return nil, errors.New(model.ServiceUnavailableMessage)
	}

	// check for duplicate email
	exists, err := s.userRepo.ExistByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New(model.GenericErrorMessage)
	}
	if exists {
		return nil, errors.New(model.EmailAlreadyExistMessage)
	}

	// hash the password before storing
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New(model.GenericErrorMessage)
	}

	// convert phone number string -> int64 (validated as numeric by the handler)
	phoneNumber, err := strconv.ParseInt(req.PhoneNumber, 10, 64) // Atoi(req.PhoneNumber)
	if err != nil {
		return nil, errors.New(model.GenericErrorMessage)
	}

	saveUserDB := &model.User{
		Email:       req.Email,
		Password:    string(hashPassword),
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: phoneNumber,
	}

	err = s.userRepo.CreateUser(ctx, saveUserDB)
	if err != nil {
		return nil, errors.New(model.GenericErrorMessage)
	}

	registerResponse := &model.AppResponse{
		Code:    model.StatusSuccess,
		Message: "Register Success",
	}

	return registerResponse, nil
}

// gives a registered user an access token if the credentials are valid
func (s *authService) Login(ctx context.Context, req *model.LoginRequest) (*model.AppResponse, error) {

	// check service hours
	if err := s.checkServiceHours(); err != nil {
		return nil, errors.New(model.ServiceUnavailableMessage)
	}

	// find user email
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New(model.InvalidCredentialMessage)
	}
	if user == nil {
		return nil, errors.New(model.InvalidCredentialMessage)
	}

	// compare user input password against stored hashed password
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	)
	if err != nil {
		return nil, errors.New(model.GenericErrorMessage)
	}

	// pass = return jwt
	// create jwt token 
	token, err := s.generateJWT(user.ID)
	if err != nil {
		return nil, errors.New(model.GenericErrorMessage)
	}

	response := &model.AppResponse{
		Code:    model.StatusSuccess,
		Message: "Login Successful",
		Data:    model.TokenData{AccessToken: token},
	}
	return response, nil
}

// register
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

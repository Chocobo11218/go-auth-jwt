package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/Chocobo11218/go-auth-jwt/app/internal/model"
	"github.com/Chocobo11218/go-auth-jwt/app/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, req *model.RegisterRequest) (model.AppResponse, error)
	Login(ctx context.Context, req *model.LoginRequest) (model.AppResponse, error)
}

type authService struct {
	//config   *config.AppConfig
	userRepo     repository.UserRepository
	redisRepo    repository.UserRepository
	jwtSecret    string
	jwtExpiresIn time.Duration
	loc          *time.Location
}

func NewAuthService(
	userRepo repository.UserRepository,
	redisRepo repository.UserRepository,
	jwtSecret string,
	jwtExpiresIn time.Duration,
	loc *time.Location,
) AuthService {
	return &authService{
		userRepo:     userRepo,
		redisRepo:    redisRepo,
		jwtSecret:    jwtSecret,
		jwtExpiresIn: jwtExpiresIn,
		loc:          loc,
	}
}

// register create a new user (ref. core-purchase-statement)
func (s *authService) Register(ctx context.Context, req *model.RegisterRequest) (model.AppResponse, error) {
	now := time.Now()

	if err := s.checkServiceHours(); err != nil {
		return model.AppResponse{}, err
	}

	exists, err := s.userRepo.ExistByEmail(ctx, req.Email)
	if err != nil {
		return model.AppResponse{}, errors.New(model.GenericErrorMessage)
	}
	if exists {
		return model.AppResponse{}, errors.New(model.EmailAlreadyExistMessage)
	}

	salt, err := generateSalt()
	if err != nil {
		return model.AppResponse{}, errors.New(model.GenericErrorMessage)
	}

	hash, err := hashPassword(req.Password, salt)
	if err != nil {
		return model.AppResponse{}, errors.New(model.GenericErrorMessage)
	}

	phoneNumber, err := strconv.Atoi(req.PhoneNumber) //ParseInt(req.PhoneNumber, 10, 64)
	if err != nil {
		return model.AppResponse{}, errors.New(model.GenericErrorMessage)
	}

	id := uuid.NewString()

	saveUserDB := &model.User{
		Id:           id,
		Email:        req.Email,
		PasswordHash: hash,
		PasswordSalt: salt,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		PhoneNumber:  int64(phoneNumber),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	err = s.userRepo.Create(ctx, saveUserDB)
	if err != nil {
		return model.AppResponse{}, errors.New(model.GenericErrorMessage)
	}

	registerResponse := model.AppResponse{
		Code:    model.StatusSuccess,
		Message: "Register Success",
	}

	return registerResponse, nil
}

func (s *authService) Login(ctx context.Context, req *model.LoginRequest) (model.AppResponse, error) {
	return model.AppResponse{}, nil
}

// register
func (s *authService) checkServiceHours() error {
	hour := time.Now().Hour()
	if hour < 6 || hour >= 23 {
		return errors.New(model.ServiceUnavailableMessage)
	}
	return nil
}

// login
// PW_SALT_BYTES = 16 -> 32-character string
func generateSalt() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func hashPassword(password, salt string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(salt+password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// verifyPassword
func verifyPassword(password, salt, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(salt+password))
}

// checkCredential -> repository.GetUserByEmail

type jwtClaims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

func (s *authService) generateJWT(user model.User) (string, error) {
	// prepare custome claims
	claims := &jwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.jwtExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID: user.Id,
		Email:  user.Email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// store session
type sessionPayload struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

// storeSession stores user session in Redis under key "session:{user_id}" with 24h TTL
func (s *authService) storeSession(ctx context.Context, user *model.User) error {
	data, err := json.Marshal(sessionPayload{ // convert the struct into a JSON byte slice
		UserID: user.Id,
		Email:  user.Email,
	})
	if err != nil {
		return err
	}
	return s.userRepo.Set(ctx, "session:"+user.Id, string(data), 24*time.Hour)
}

/*
{
  "user_id": 1,
  "email": "customer@example.com",
  "exp": 1234567890
}
*/

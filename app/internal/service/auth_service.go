package service

import (
	"context"
	"errors"
	"time"

	"github.com/Chocobo11218/go-auth-jwt/app/internal/model"
	"github.com/Chocobo11218/go-auth-jwt/app/internal/repository"
	"github.com/Chocobo11218/go-auth-jwt/app/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// AuthService is the port (interface) that the HTTP handler depends on.

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
		logger.Error(ctx, "Register Service: Service hours restricted",
			zap.Error(err),
		)
		return model.AppResponse{}, errors.New(model.ServiceUnavailableMessage)
	}

	// if err := validatePasswordStrength(req.Password); err != nil {
	//     return nil, errors.New(model.GenericErrorMessage)
	// }

	// check for duplicate email
	exists, err := s.userRepo.ExistByEmail(ctx, req.Email)

	if err != nil {
		logger.Error(ctx, "Register Service: Failed to check email existence",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return model.AppResponse{}, errors.New(model.GenericErrorMessage)
	}
	if exists { // err == nil
		logger.Error(ctx, "Register Service: Email already exists",
			zap.String("email", req.Email),
		)
		return model.AppResponse{}, errors.New(model.EmailAlreadyExistMessage)
	}

	// hash the password before storing
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error(ctx, "Register Service: Failed to hash password",
			zap.Error(err),
		)
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
		logger.Error(ctx, "Register Service: Failed to create user",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		// if err.Error() == model.EmailAlreadyExistMessage {
		//     return model.AppResponse{}, err
		// }
		return model.AppResponse{}, errors.New(model.GenericErrorMessage)
	}

	registerResponse := model.AppResponse{
		Code:    model.StatusSuccess,
		Message: "Register Success",
	}

	return registerResponse, nil
}

// gives a registered user an access token if the credentials are valid
var dummyHash, _ = bcrypt.GenerateFromPassword([]byte("dummy-password"), bcrypt.DefaultCost)

func (s *authService) Login(ctx context.Context, req *model.LoginRequest) (model.AppResponse, error) {

	// check service hours
	if err := s.checkServiceHours(); err != nil {
		logger.Error(ctx, "Login Service: Service hours restricted",
			zap.Error(err),
		)
		return model.AppResponse{}, errors.New(model.ServiceUnavailableMessage)
	}

	// find user email
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)

	if err != nil {
		logger.Error(ctx, "Login Service: User not found",
			zap.Error(err),
		)
		return model.AppResponse{}, errors.New(model.GenericErrorMessage)
	}
	// no user found
	if user == nil {
		logger.Error(ctx, "Login Service: User look up failed")
		// run bcrypt for constant time
		bcrypt.CompareHashAndPassword(dummyHash, []byte(req.Password))

		return model.AppResponse{}, errors.New(model.InvalidCredentialMessage)
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
		return model.AppResponse{}, errors.New(model.InvalidCredentialMessage)
	}

	// pass = return jwt
	// create jwt token
	token, err := s.generateJWT(user.ID)
	if err != nil {
		logger.Error(ctx, "Login Service: jwt generation failed",
			zap.Error(err),
		)
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
	hour := time.Now().In(time.Local).Hour() // loc *time.Location
	if hour < 6 || hour >= 23 {
		return errors.New(model.ServiceUnavailableMessage)
	}
	return nil
}

func (s *authService) generateJWT(userID uint) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["jti"] = uuid.New().String() // unique ID, needed for revocation later
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(s.jwtTTL).Unix()

	t, err := token.SignedString([]byte(s.jwtSecretKey))
	if err != nil {
		return "", err
	}
	return t, nil
}

// func validatePasswordStrength(password string) error {
//     var hasUpper, hasLower, hasDigit, hasSpecial bool
//     for _, c := range password {
//         switch {
//         case unicode.IsUpper(c): hasUpper = true
//         case unicode.IsLower(c): hasLower = true
//         case unicode.IsDigit(c): hasDigit = true
//         case unicode.IsPunct(c) || unicode.IsSymbol(c): hasSpecial = true
//         }
//     }
//     if !hasUpper || !hasLower || !hasDigit || !hasSpecial {
//         return errors.New("password must contain upper, lower, digit, and special character")
//     }
//     return nil
// }

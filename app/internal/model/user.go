package model

import "time"

// db entity
type User struct {
	Id           string     `gorm:"column:id"`
	Email        string     `gorm:"column:email"`
	PasswordHash string     `gorm:"column:password_hash"`
	PasswordSalt string     `gorm:"column:password_salt"`
	FirstName    string     `gorm:"column:first_name"`
	LastName     string     `gorm:"column:last_name"`
	PhoneNumber  int64      `gorm:"column:phone_number"`
	CreatedAt    time.Time  `gorm:"column:created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at"`
}

// request
type RegisterRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required,numeric"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// response
type AppResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type TokenData struct {
	AccessToken string `json:"access_token"`
}

// error
const (
	// code
	StatusSuccess            = 1000
	StatusEmailAlreadyExist  = 4001
	StatusInvalidCredential  = 4002
	StatusGenericError       = 5000
	StatusServiceUnavailable = 5001
	StatusBadRequest         = 400

	// message
	ServiceUnavailableMessage = "Service is available between 06:00 and 23:00"
	EmailAlreadyExistMessage  = "Email already exist"
	InvalidCredentialMessage  = "Invalid email or password"
	GenericErrorMessage = "Generic error"
)

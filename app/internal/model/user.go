package model

import "gorm.io/gorm"

// db entity
type User struct {
	gorm.Model         // ID, Created_at, Updated_at, Deleted_at
	Email       string `gorm:"unique"`
	Password    string
	FirstName   string
	LastName    string
	PhoneNumber int64
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
	GenericErrorMessage       = "Generic error"
)

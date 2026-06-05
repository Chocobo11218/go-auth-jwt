package model

import (
	"time"

	"gorm.io/gorm"
)

// db entity
type User struct {
	ID          uint           `gorm:"primaryKey;autoIncrement;not null"`
	Email       string         `gorm:"unique;size:255;not null"`
	Password    string         `gorm:"size:255;not null"`
	FirstName   string         `gorm:"size:100;not null"`
	LastName    string         `gorm:"size:100;not null"`
	PhoneNumber string         `gorm:"size:20;not null"`
	Created_at  time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP"`
	Updated_at  time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP"`
	Deleted_at  gorm.DeletedAt `gorm:"index"`
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

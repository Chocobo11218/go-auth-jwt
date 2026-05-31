package model

import "time"

// db entity
type User struct {
	Id           string     `gorm:"column:id"`
	Email        string     `gorm:"column:email"`
	PasswordHash string     `gorm:"column:password_hash"`
	PasswordSalt string     `gorm:"column:password_hash"`
	FirstName    string     `gorm:"column:password_hash"`
	LastName     string     `gorm:"column:password_hash"`
	PhoneNumber  int64      `gorm:"column:password_hash"`
	CreatedAt    time.Time  `gorm:"column:password_hash"`
	UpdatedAt    time.Time  `gorm:"column:password_hash"`
	DeletedAt    *time.Time `gorm:"column:password_hash"`
}

// request
type RegisterRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
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
	EmailAlreadyExistMessage = "Email already exist"
	ServiceUnavailableMessage = "Service is available between 06:00 and 23:00"

	StatusSuccess = 1000
	StatusEmailAlreadyExist = 4001
	StatusInvalidCredential = 4002
	StatusGenericError = 5000
	StatusServiceUnavailable = 5001
	StatusBadRequest = 400
)

type AppError struct {
	Code int
	Message string
}
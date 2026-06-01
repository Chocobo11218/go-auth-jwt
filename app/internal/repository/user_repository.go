package repository

import (
	"context"
	"errors"

	"github.com/Chocobo11218/go-auth-jwt/app/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	ExistByEmail(ctx context.Context, email string) (bool, error)
	Create(ctx context.Context, req *model.User) error
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// 
func (r *userRepository) ExistByEmail(ctx context.Context, email string) (bool, error) {

	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("email = ? AND deleted_at IS NULL", email).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	
	return count > 0, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	
	var user model.User
	err := r.db.WithContext(ctx).
		Where("email = ? AND deleted_at IS NULL", email).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound){
			return nil, nil
		}
		return nil, err
	}
	
	return &user, nil
}

// inserts a new user into the database
func (r *userRepository) Create(ctx context.Context, req *model.User) error {
	//req.Id = uuid.New().String()
	return r.db.WithContext(ctx).Create(req).Error
}

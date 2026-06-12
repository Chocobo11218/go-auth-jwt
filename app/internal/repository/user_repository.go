package repository

import (
	"context"
	"errors"

	"github.com/Chocobo11218/go-auth-jwt/app/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	ExistByEmail(ctx context.Context, email string) (bool, error)
	CreateUser(ctx context.Context, req *model.User) error
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByID(ctx context.Context, id uint) (*model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) ExistByEmail(ctx context.Context, email string) (bool, error) {
	var exists int
	err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Select("1").
		Where("email = ? AND deleted_at IS NULL", email).
		Limit(1).
		Scan(&exists).Error
	// SELECT 1 FROM users WHERE ... LIMIT 1
	return exists == 1, err
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {

	var user model.User
	err := r.db.WithContext(ctx).
		Where("email = ? AND deleted_at IS NULL", email).
		First(&user).Error
	// SELECT * FROM `users` WHERE (email = 'test@example.com' AND deleted_at IS NULL) AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// inserts a new user into the database
func (r *userRepository) CreateUser(ctx context.Context, user *model.User) error {
	
	return r.db.WithContext(ctx).Create(user).Error
}

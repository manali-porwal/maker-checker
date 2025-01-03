package user

import (
	"context"
	"errors"
	"maker-checker/internal/models"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, user models.User) (*models.User, error) {
	user.CreatedAt = time.Now()

	if err := r.db.
		WithContext(ctx).
		Clauses(clause.Returning{}).
		Create(&user).
		Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetByUserName(ctx context.Context, username string) (*models.User, error) {
	var user models.User

	err := r.db.
		WithContext(ctx).
		Where("username = ?", username).
		First(&user).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

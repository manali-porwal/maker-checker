package message

import (
	"context"
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

func (r *Repository) Create(ctx context.Context, message models.Message) (*models.Message, error) {
	message.CreatedAt = time.Now()

	if err := r.db.
		WithContext(ctx).
		Clauses(clause.Returning{}).
		Create(&message).
		Error; err != nil {
		return nil, err
	}

	return &message, nil
}

func (r *Repository) Get(ctx context.Context, id uint) (*models.Message, error) {
	var message models.Message

	if err := r.db.
		WithContext(ctx).
		Clauses(clause.Returning{}).
		First(&message, id).
		Error; err != nil {
		return nil, err
	}

	return &message, nil
}

func (r *Repository) GetApproval(ctx context.Context, messageID uint, approverID uint) (*models.MessageApproval, error) {
	var messageApproval models.MessageApproval

	if err := r.db.
		WithContext(ctx).
		Where("message_id = ? AND approver_id = ?", messageID, approverID).
		First(&messageApproval).Error; err != nil {
		return nil, err
	}

	return &messageApproval, nil
}

func (r *Repository) UpdateStatus(ctx context.Context, message models.Message) (*models.Message, error) {
	if err := r.db.
		WithContext(ctx).
		Model(&message).
		Update("status", message.Status).
		Error; err != nil {
		return nil, err
	}

	return &message, nil
}

func (r *Repository) UpdateApproval(ctx context.Context, message models.Message, approval models.MessageApproval,
) (*models.Message, *models.MessageApproval, error) {
	tx := r.db.Begin()

	if err := r.db.
		WithContext(ctx).
		Clauses(clause.Returning{}).
		Where("id=?", message.ID).
		Updates(&message).
		Error; err != nil {
		return nil, nil, err
	}

	approval.CreatedAt = time.Now()
	if err := r.db.
		WithContext(ctx).
		Clauses(clause.Returning{}).
		Create(&approval).
		Error; err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	tx.Commit()
	return &message, &approval, nil
}

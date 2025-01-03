package message

import (
	"context"
	"fmt"
	"log"
	"maker-checker/config"
	"maker-checker/internal/models"
)

type MessageRepository interface {
	Create(ctx context.Context, message models.Message) (*models.Message, error)
	Get(ctx context.Context, id uint) (*models.Message, error)
	UpdateApproval(ctx context.Context, message models.Message, approval models.MessageApproval,
	) (*models.Message, *models.MessageApproval, error)
	UpdateStatus(ctx context.Context, message models.Message) (*models.Message, error)
	GetApproval(ctx context.Context, messageID uint, approverID uint) (*models.MessageApproval, error)
}

type MessageService struct {
	messageRepo          MessageRepository
	numRequiredApprovals int
}

func New(messageRepo MessageRepository, cfg *config.AppConfig) *MessageService {
	return &MessageService{
		messageRepo:          messageRepo,
		numRequiredApprovals: cfg.Checker.NumRequiredApprovals,
	}
}

func (s *MessageService) CreateMessage(ctx context.Context, msg *models.Message) (*models.Message, error) {
	// Business logic for message creation
	msg.Status = models.Pending
	msg.RequiredApprovals = s.numRequiredApprovals // this is a configurable field

	msg, err := s.messageRepo.Create(ctx, *msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (s *MessageService) Get(ctx context.Context, id uint) (*models.Message, error) {
	msg, err := s.messageRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (s *MessageService) ProcessApproval(ctx context.Context, approval *models.MessageApproval) (*models.Message, error) {
	msg, err := s.messageRepo.Get(ctx, approval.MessageID)
	if err != nil {
		return nil, fmt.Errorf("message not found, error: %v", err)
	}

	if msg.Status != models.Pending {
		return nil, fmt.Errorf("message is not in pending state")
	}

	//check if user has already approved or rejected this message
	_, err = s.messageRepo.GetApproval(ctx, approval.MessageID, approval.ApproverID)
	if err == nil {
		return nil, fmt.Errorf("user has already processed this message")
	}

	if approval.Approved {
		msg.ApprovalCount++
		if msg.ApprovalCount == s.numRequiredApprovals {
			msg.Status = models.Approved
			// Simulate sending message
			go s.SendMessage(ctx, msg.ID)
		}
	} else {
		msg.RejectionCount++
		msg.Status = models.Rejected
	}

	msg, _, err = s.messageRepo.UpdateApproval(ctx, *msg, *approval)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (s *MessageService) SendMessage(ctx context.Context, msgID uint) (*models.Message, error) {
	msg, err := s.messageRepo.Get(ctx, msgID)
	if err != nil {
		return nil, fmt.Errorf("message not found, error: %v", err)
	}

	if msg.Status != models.Approved {
		return nil, fmt.Errorf("message is not in approved")
	}

	msg.Status = models.Delivered
	log.Printf("Sending the message ID %d created by user ID %s to recipient %s", msg.ID, msg.CreatedBy, msg.Recipient)
	return s.messageRepo.UpdateStatus(ctx, *msg)
}

package models

import (
	"time"
)

type MessageStatus string

const (
	Pending   MessageStatus = "pending"
	Approved  MessageStatus = "approved"
	Rejected  MessageStatus = "rejected"
	Delivered MessageStatus = "delivered"
)

type Message struct {
	ID                uint          `json:"id"`
	Content           string        `json:"content"`
	Recipient         string        `json:"recipient"`
	Status            MessageStatus `json:"status"`
	CreatedBy         uint          `json:"created_by"`
	CreatedAt         time.Time     `json:"created_at" gorm:"created_at"`
	RequiredApprovals int           `json:"required_approvals"`
	ApprovalCount     int           `json:"approval_count"`
	RejectionCount    int           `json:"rejection_count"`
}

type MessageApproval struct {
	MessageID  uint      `json:"message_id" gorm:"message_id"`
	Approved   bool      `json:"approved" gorm:"approved"`
	ApproverID uint      `json:"approver_id" gorm:"approver_id"`
	Comment    string    `json:"comment" gorm:"comment"`
	CreatedAt  time.Time `json:"created_at" gorm:"created_at"`
}

package dtos

type MessageStatus string

// CreateMessageRequest create a new message
type CreateMessageRequest struct {
	Content   string `binding:"required" json:"content"`
	Recipient string `binding:"required" json:"recipient"`
}

// ApprovalRequest represents an approval/rejection request
type ApprovalRequest struct {
	Approved *bool  `binding:"required" json:"approved"`
	Comment  string `binding:"omitempty" json:"comment"`
}

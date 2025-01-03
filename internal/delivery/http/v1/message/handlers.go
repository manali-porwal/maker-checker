package message

import (
	"context"
	crtlModels "maker-checker/internal/dtos"
	"maker-checker/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MessageService interface {
	CreateMessage(ctx context.Context, msg *models.Message) (*models.Message, error)
	Get(ctx context.Context, id uint) (*models.Message, error)
	ProcessApproval(ctx context.Context, approval *models.MessageApproval) (*models.Message, error)
}

type MessageHandler struct {
	messageService MessageService
}

func New(messageService MessageService) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
	}
}

func (h *MessageHandler) Create(c *gin.Context) {
	req := crtlModels.CreateMessageRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})

		return
	}

	maker_id, found := c.Get("user_id")
	if !found {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID missing"})
		return
	}
	message, err := h.messageService.CreateMessage(c,
		&models.Message{
			Content:   req.Content,
			Recipient: req.Recipient,
			CreatedBy: maker_id.(uint),
		})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusCreated, message)
}

func (h *MessageHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})

		return
	}

	message, err := h.messageService.Get(c, uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, message)
}

func (h *MessageHandler) Approval(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})

		return
	}

	req := crtlModels.ApprovalRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})

		return
	}

	approver_id, found := c.Get("user_id")
	if !found {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID missing"})
		return
	}

	message, err := h.messageService.ProcessApproval(c,
		&models.MessageApproval{
			MessageID:  uint(id),
			Approved:   *req.Approved,
			ApproverID: approver_id.(uint),
			Comment:    req.Comment,
		})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, message)
}

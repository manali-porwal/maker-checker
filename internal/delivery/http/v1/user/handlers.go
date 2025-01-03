package user

import (
	"context"
	"maker-checker/internal/dtos"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	Create(ctx context.Context, request *dtos.CreateUserRequest) (*dtos.CreateUserResponse, error)
	Login(ctx context.Context, request *dtos.LoginRequest) (*dtos.LoginResponse, error)
}

type UserHandler struct {
	userService UserService
}

func New(userService UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Create(c *gin.Context) {
	request := dtos.CreateUserRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})

		return
	}

	response, err := h.userService.Create(c, &request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": response})
}

func (h *UserHandler) Login(c *gin.Context) {
	request := dtos.LoginRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})

		return
	}

	response, err := h.userService.Login(c, &request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": response})
}

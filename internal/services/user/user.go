package user

import (
	"context"
	"fmt"
	"maker-checker/config"
	"maker-checker/internal/dtos"
	"maker-checker/internal/models"
	"maker-checker/pkg/auth"
	"time"

	"github.com/golang-jwt/jwt"
)

type UserRepository interface {
	Create(ctx context.Context, user models.User) (*models.User, error)
	GetByUserName(ctx context.Context, username string) (*models.User, error)
}

type UserService struct {
	userRepo UserRepository
	jwt      *auth.JWT
}

func New(userRepo UserRepository, cfg *config.AppConfig) *UserService {
	return &UserService{
		userRepo: userRepo,
		jwt:      auth.New(cfg),
	}
}

func (s *UserService) Create(ctx context.Context, request *dtos.CreateUserRequest) (*dtos.CreateUserResponse, error) {
	user, err := s.userRepo.GetByUserName(ctx, request.UserName)
	if err != nil {
		return nil, err
	}

	if user != nil {
		return nil, fmt.Errorf("user account already exists")
	}

	hashedPassword, err := auth.HashPassword(request.Password)
	if err != nil {
		return nil, err
	}

	userResponse, err := s.userRepo.Create(ctx, models.User{
		Username:       request.UserName,
		HashedPassword: hashedPassword,
		Role:           request.Role,
	})
	if err != nil {
		return nil, err
	}

	loginResponse, err := s.generateAccessToken(ctx, userResponse)
	if err != nil {
		return nil, err
	}

	response := dtos.CreateUserResponse{
		UserID:      userResponse.ID,
		AccessToken: loginResponse.AccessToken,
		ExpiredAt:   loginResponse.ExpiredAt,
	}
	return &response, nil
}

func (s *UserService) Login(ctx context.Context, request *dtos.LoginRequest) (*dtos.LoginResponse, error) {
	user, err := s.userRepo.GetByUserName(ctx, request.Username)
	if err != nil {
		return nil, err
	}

	verified, err := auth.VerifyPassword(user.HashedPassword, request.Password)
	if err != nil {
		return nil, err
	}

	if !verified {
		return nil, fmt.Errorf("invalid username or password")
	}

	return s.generateAccessToken(ctx, user)
}

func (s *UserService) generateAccessToken(ctx context.Context, user *models.User) (*dtos.LoginResponse, error) {
	expiresAt := time.Now().Add(time.Hour * time.Duration(s.jwt.ExpiryHours))

	claims := auth.PayloadToken{
		UserID:   user.ID,
		UserName: user.Username,
		Role:     user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
		},
	}

	// Generate encoded token and send it as response.
	jwtToken, err := s.jwt.GenerateToken(ctx, claims)
	if err != nil {
		return nil, err
	}

	return &dtos.LoginResponse{
		AccessToken: jwtToken,
		ExpiredAt:   expiresAt,
	}, nil
}

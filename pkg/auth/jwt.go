package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"maker-checker/config"
	"strings"

	"github.com/golang-jwt/jwt"
)

type (
	PayloadToken struct {
		UserID   uint   `json:"user_id"`
		UserName string `json:"user_name"`
		Role     string `json:"role"`
		jwt.StandardClaims
	}

	Data struct {
		UserID   uint   `json:"user_id"`
		UserName string `json:"user_name"`
		Role     string `json:"role"`
	}
)

var (
	errMissingJwtToken = errors.New("Missing JWT token")
	errInvalidJwtToken = errors.New("Invalid JWT token")
)

type JWT struct {
	Secret      string
	ExpiryHours int
}

func New(cfg *config.AppConfig) *JWT {
	return &JWT{Secret: cfg.JWT.Secret, ExpiryHours: cfg.JWT.ExpiryHours}
}

func (j *JWT) GenerateToken(ctx context.Context, request PayloadToken) (token string, err error) {

	return jwt.NewWithClaims(jwt.SigningMethodHS256, request).SignedString([]byte(j.Secret))
}

func (j *JWT) extractClaims(tokenStr string) (*PayloadToken, error) {
	hmacSecret := []byte(j.Secret)
	token, err := jwt.ParseWithClaims(tokenStr, &PayloadToken{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errInvalidJwtToken
		}
		return hmacSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*PayloadToken); ok && token.Valid {
		return claims, nil
	} else {
		log.Println("Invalid JWT Token", ok, token.Valid)
		return nil, errInvalidJwtToken
	}
}

func (j *JWT) ValidateToken(ctx context.Context, token string) (*Data, error) {
	fmt.Println("1 ", j.Secret)
	tokenWithoutBearer := strings.Replace(token, "Bearer ", "", 1)
	if len(tokenWithoutBearer) == 0 {
		return nil, errMissingJwtToken
	}

	fmt.Println("2 ", tokenWithoutBearer)

	//newToken, err := jwt.ParseWithClaims(tokenWithoutBearer, &PayloadToken{}, func(token *jwt.Token) (interface{}, error) {
	claims, err := j.extractClaims(tokenWithoutBearer)
	if err != nil {
		return nil, err
	}
	return &Data{
		UserID:   claims.UserID,
		UserName: claims.UserName,
		Role:     claims.Role,
	}, nil
}

// 	// Store the token claims in the request context for later use
// 	ctx.Set(constant.AuthCredentialKey, newToken.Claims.(*PayloadToken))

// 	return nil
// }

// func (j *JWT) RevokeToken(ctx context.Context, token string, expiration time.Duration) error {
// 	return j.redis.Set(ctx, token, constant.TokenRevoked, expiration)
// }

// func (j *JWT) IsTokenRevoked(ctx context.Context, token string) bool {
// 	revoked, err := j.redis.Get(ctx, token)
// 	if err != nil {
// 		return false
// 	}

// 	return revoked == constant.TokenRevoked
// }

// func NewTokenInformation(ctx context.Context) (*PayloadToken, error) {
// 	tokenInformation, ok := ctx.Get(constant.AuthCredentialKey).(*PayloadToken)
// 	if !ok {
// 		return tokenInformation, apperror.Unauthorized(apperror.ErrFailedGetTokenInformation)
// 	}

// 	return tokenInformation, nil
// }

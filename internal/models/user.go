package models

import "time"

const (
	RoleMaker   = "MAKER"
	RoleChecker = "CHECKER"
	RoleAdmin   = "ADMIN"
)

type User struct {
	ID             uint      `json:"id"`
	Username       string    `json:"username"`
	HashedPassword string    `json:"-"` // Never send password in response
	Role           string    `json:"role"`
	CreatedAt      time.Time `json:"created_at" gorm:"created_at"`
}

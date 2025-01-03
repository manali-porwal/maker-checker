package middleware

import (
	"maker-checker/pkg/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Middleware struct {
	jwt *auth.JWT
}

func New(jwt *auth.JWT) *Middleware {
	return &Middleware{jwt: jwt}
}

// Middleware function to validate JWT token
func (m *Middleware) AuthMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader(string("Authorization"))
		if tokenString == "" {
			ctx.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"error": "token missing"},
			)
			return
		}

		data, err := m.jwt.ValidateToken(ctx, tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"error": err.Error()},
			)
			return
		}

		// Check if user role is allowed
		roleAllowed := false
		for _, role := range allowedRoles {
			if data.Role == role {
				roleAllowed = true
				break
			}
		}

		if !roleAllowed {
			ctx.AbortWithStatusJSON(
				http.StatusForbidden,
				gin.H{"error": "forbidden"},
			)
			return
		}

		ctx.Set("user_id", data.UserID)
		ctx.Set("user_name", data.UserName)
		ctx.Set("role", data.Role)

		ctx.Next()
	}
}

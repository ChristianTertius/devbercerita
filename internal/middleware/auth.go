package middleware

import (
	"ChristianTertius/devbercerita/pkg/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

func OptionalAuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.Request.Header.Get("Authorization")
		header = strings.TrimSpace(header)
		if header == "" {
			c.Next()
			return
		}
		userID, username, err := jwt.ValidateToken(header, secretKey, true)
		if err != nil {
			c.Next()
			return
		}
		c.Set("userID", userID)
		c.Set("username", username)
		c.Next()
	}
}

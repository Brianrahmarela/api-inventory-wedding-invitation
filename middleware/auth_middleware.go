package middleware

import (
	"api-go-invitation/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware memverifikasi token JWT dan menyimpan userId & role ke dalam context Gin
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization is required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, gin.H{"error": "Authorization header must be Bearer Token"})
			c.Abort()
			return
		}

		// Validasi token & ambil userId serta role
		userId, role, err := utils.ValidateToken(parts[1])
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Simpan userId dan role ke dalam context untuk bisa digunakan di handler
		c.Set("userId", userId)
		c.Set("role", role)

		c.Next()
	}
}

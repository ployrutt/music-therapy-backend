package middleware

import (
	"net/http"
	"strings"

	"project-backend/helpers"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware ตรวจสอบ JWT และตรวจสอบสิทธิ์ตาม Roles
func AuthMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. ดึง Token จาก Header "Authorization: Bearer <token>"
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 2. ตรวจสอบ Token
		claims, err := helpers.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// 3. ตรวจสอบ Role (Authorization Check)
		if len(allowedRoles) > 0 {
			isAllowed := false
			for _, role := range allowedRoles {
				if claims.RoleName == role {
					isAllowed = true
					break
				}
			}
			if !isAllowed {
				// 403 Forbidden
				c.JSON(http.StatusForbidden, gin.H{"error": "Access denied for this role"})
				c.Abort()
				return
			}
		}

		// 4. บันทึกข้อมูล UserID และ RoleName ลงใน Context
		c.Set("user_id", claims.UserID)
		c.Set("role_name", claims.RoleName)

		c.Next() // ไปยัง Controller ถัดไป
	}
}

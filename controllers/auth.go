package controllers

import (
	"net/http"
	"time"

	"project-backend/helpers"
	"project-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RegisterInput struct {
	FirstName   string    `json:"first_name" binding:"required"`
	LastName    string    `json:"last_name" binding:"required"`
	Email       string    `json:"email" binding:"required,email"`
	Password    string    `json:"password" binding:"required,min=8"`
	PhoneNumber string    `json:"phone_number" binding:"required"`
	DateOfBirth time.Time `json:"date_of_birth" binding:"required"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type ForgotPasswordInput struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordInput struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

func Register(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input RegisterInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var existingUser models.User
		if err := db.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email already registered."})
			return
		}

		hashedPassword, err := helpers.HashPassword(input.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password."})
			return
		}

		var memberRole models.Role
		if err := db.Where("role_name = ?", "member").First(&memberRole).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Member role not found."})
			return
		}

		dateOnly := input.DateOfBirth.Truncate(24 * time.Hour)
		user := models.User{
			FirstName:   input.FirstName,
			LastName:    input.LastName,
			DateOfBirth: dateOnly,
			Email:       input.Email,
			Password:    hashedPassword,
			PhoneNumber: input.PhoneNumber,
			RoleID:      memberRole.ID,
		}
		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user."})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Registration successful.", "user_id": user.ID})
	}
}

func Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input LoginInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var user models.User
		if err := db.Where("LOWER(email) = LOWER(?)", input.Email).Preload("Role").First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		if !helpers.CheckPasswordHash(input.Password, user.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
		roleName := user.Role.RoleName
		token, err := helpers.GenerateToken(user.ID, roleName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Login successful",
			"token":   token,
			"role":    roleName,
		})
	}
}

// ForgotPassword: ตรวจสอบอีเมลและสร้างรหัสรีเซ็ต
func ForgotPassword(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Email string `json:"email" binding:"required,email"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "กรุณากรอกอีเมลให้ถูกต้อง"})
			return
		}

		var user models.User
		if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
			// ตอบ OK เพื่อไม่ให้ Hacker ทราบว่ามีอีเมลนี้ในระบบหรือไม่ (Security Best Practice)
			c.JSON(http.StatusOK, gin.H{"message": "หากพบอีเมลในระบบ ระบบจะส่งรหัสรีเซ็ตไปให้ท่าน"})
			return
		}

		token, err := helpers.GenerateResetToken(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถสร้างรหัสรีเซ็ตได้"})
			return
		}

		// ในอนาคต: เขียนฟังก์ชันส่งอีเมลแนบ token ตรงนี้
		c.JSON(http.StatusOK, gin.H{
			"message": "สร้างรหัสรีเซ็ตสำเร็จ (ในระบบจริงจะส่งเข้าอีเมล)",
			"token":   token, // ส่งออกไปเพื่อให้คุณทดสอบผ่าน Postman ได้ก่อน
		})
	}
}

func ResetPassword(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Token       string `json:"token" binding:"required"`
			NewPassword string `json:"new_password" binding:"required,min=8"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ข้อมูลไม่ถูกต้องหรือรหัสผ่านสั้นเกินไป"})
			return
		}

		// 1. ตรวจสอบความถูกต้องของ Token
		claims, err := helpers.ValidateToken(input.Token)
		if err != nil || claims.RoleName != "password_reset" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "รหัสรีเซ็ตไม่ถูกต้องหรือหมดอายุแล้ว"})
			return
		}

		// 2. Hash รหัสผ่านใหม่ (ใช้ตัวเดียวกับที่ใช้ใน Register)
		hashedPassword, err := helpers.HashPassword(input.NewPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "เกิดข้อผิดพลาดในการตั้งรหัสผ่าน"})
			return
		}

		// 3. อัปเดตลง Database
		if err := db.Model(&models.User{}).Where("id = ?", claims.UserID).Update("password", hashedPassword).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถอัปเดตรหัสผ่านได้"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "เปลี่ยนรหัสผ่านใหม่สำเร็จแล้ว"})
	}
}

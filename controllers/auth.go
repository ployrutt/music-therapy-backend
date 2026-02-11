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

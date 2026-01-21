package controllers

import (
	"net/http"

	"project-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ListAllUsers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []models.User

		// 1. เพิ่ม "date_of_birth" (หรือชื่อคอลัมน์จริงใน DB) เข้าไปใน Select
		if err := db.
			Select("id", "firstname", "lastname", "email", "phone", "role_id", "created_at", "date_of_birth").
			Preload("Role").
			Find(&users).Error; err != nil {
			// ... error handling ...
		}

		type UserResponse struct {
			ID          uint   `json:"id"`
			FirstName   string `json:"first_name"`
			LastName    string `json:"last_name"`
			Email       string `json:"email"`
			PhoneNumber string `json:"phone_number"`
			DateOfBirth string `json:"date_of_birth"` // เปลี่ยนเป็น string เพื่อ format
			RoleName    string `json:"role_name"`
		}

		var responseData []UserResponse
		for _, user := range users {
			roleName := ""
			if user.Role != nil {
				roleName = user.Role.RoleName
			}

			responseData = append(responseData, UserResponse{
				ID:          user.ID,
				FirstName:   user.FirstName,
				LastName:    user.LastName,
				Email:       user.Email,
				PhoneNumber: user.PhoneNumber,
				// 2. Format วันที่ให้เป็น string
				DateOfBirth: user.DateOfBirth.Format("2006-01-02"),
				RoleName:    roleName,
			})
		}

		c.JSON(http.StatusOK, responseData)
	}
}

func GetProfile(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		userID, exists := c.Get("user_id")
		if !exists {

			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
			return
		}

		var user models.User

		if err := db.
			Preload("Role").
			Omit("password", "deleted_at").
			First(&user, userID).Error; err != nil {

			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve profile data"})
			return
		}

		type ProfileResponse struct {
			ID          uint   `json:"id"`
			FirstName   string `json:"first_name"`
			LastName    string `json:"last_name"`
			Email       string `json:"email"`
			PhoneNumber string `json:"phone_number"`
			DateOfBirth string `json:"date_of_birth"`
			RoleName    string `json:"role_name"`
			Profile     string `json:"profile"`
		}

		roleName := ""
		if user.Role != nil {
			roleName = user.Role.RoleName
		}

		response := ProfileResponse{
			ID:          user.ID,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			Profile:     user.Profile,
			DateOfBirth: user.DateOfBirth.Format("2006-01-02"),
			RoleName:    roleName,
		}

		c.JSON(http.StatusOK, response)
	}
}

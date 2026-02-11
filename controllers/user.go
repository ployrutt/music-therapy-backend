package controllers

import (
	"net/http"
	"time"

	"project-backend/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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

// AdminCreateUser - Admin เพิ่มสมาชิกใหม่เอง (เช่น เพิ่ม Staff หรือ Admin คนอื่น)
func AdminCreateUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			FirstName   string    `json:"first_name" binding:"required"`
			LastName    string    `json:"last_name" binding:"required"`
			Email       string    `json:"email" binding:"required,email"`
			Password    string    `json:"password" binding:"required,min=6"`
			PhoneNumber string    `json:"phone_number" binding:"required"`
			RoleID      uint      `json:"role_id" binding:"required"`
			DateOfBirth time.Time `json:"date_of_birth"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ข้อมูลไม่ถูกต้องหรือใส่ไม่ครบ"})
			return
		}

		// 1. Hash Password
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

		// 2. สร้าง Model User
		newUser := models.User{
			FirstName:   input.FirstName,
			LastName:    input.LastName,
			Email:       input.Email,
			Password:    string(hashedPassword),
			PhoneNumber: input.PhoneNumber,
			RoleID:      input.RoleID,
			DateOfBirth: input.DateOfBirth,
		}

		// 3. บันทึกลง DB
		if err := db.Create(&newUser).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถเพิ่มสมาชิกได้ (Email อาจซ้ำ)"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "เพิ่มสมาชิกสำเร็จ", "user_id": newUser.ID})
	}
}

// AdminDeleteUser - Admin ลบสมาชิก
func AdminDeleteUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		// ใช้ Unscoped() หากต้องการลบออกจาก DB จริงๆ
		// หรือไม่ใช้เพื่อทำ Soft Delete (ตามที่คุณตั้งค่า DeletedAt ไว้ใน Model)
		if err := db.Delete(&models.User{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถลบสมาชิกได้"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "ลบสมาชิกเรียบร้อยแล้ว"})
	}
}

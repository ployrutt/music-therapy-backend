package controllers

import (
	"net/http"
	"os"

	"project-backend/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UpdateProfileRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Profile     string `json:"profile"`
	Password    string `json:"password"`
	OldPassword string `json:"old_password"`
}

func UpdateProfile(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("user_id").(uint)
		var req UpdateProfileRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "รูปแบบข้อมูลไม่ถูกต้อง"})
			return
		}

		var user models.User
		if err := db.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบผู้ใช้งาน"})
			return
		}

		updates := map[string]interface{}{}

		// --- ส่วนการจัดการข้อมูลทั่วไป ---
		if req.FirstName != "" {
			updates["firstname"] = req.FirstName
		}
		if req.LastName != "" {
			updates["lastname"] = req.LastName
		}
		if req.PhoneNumber != "" {
			updates["phone"] = req.PhoneNumber
		}
		if req.Profile != "" {
			updates["profile"] = req.Profile
		}

		if req.Password != "" {

			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "การเข้ารหัสผิดพลาด"})
				return
			}
			updates["password"] = string(hashedPassword)
		}

		if len(updates) == 0 {
			c.JSON(http.StatusOK, gin.H{"message": "ไม่มีการเปลี่ยนแปลงข้อมูล"})
			return
		}

		if err := db.Model(&user).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "บันทึกข้อมูลไม่สำเร็จ", "detail": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "แก้ไขโปรไฟล์และรหัสผ่านเรียบร้อยแล้ว"})
	}
}
func DeleteProfileImage(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		uid, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: No user_id in context"})
			return
		}

		userID, ok := uid.(uint)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error: Invalid user_id type"})
			return
		}

		var user models.User
		if err := db.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		if user.Profile != "" {
			os.Remove(user.Profile)
			db.Model(&user).Update("profile", "")
		}

		c.JSON(http.StatusOK, gin.H{"message": "Profile image deleted successfully"})
	}
}

package seeds

import (
	"errors"
	"project-backend/models"

	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB, user *models.User) error {
	if user == nil {
		return errors.New("user object cannot be nil")
	}

	// db.Create() จะทำการบันทึกข้อมูล
	// ในการ Seeding เราอาจใช้ FirstOrCreate เพื่อหลีกเลี่ยงการสร้างซ้ำ

	// ถ้าใช้ Create (สร้างใหม่ทุกครั้ง):
	result := db.Create(user)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("failed to create seed user: no rows affected")
	}

	return nil
}

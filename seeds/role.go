package seeds

import (
	"log"

	"project-backend/models"

	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) error {
	roles := []string{"admin", "member", "guest"}

	for _, roleName := range roles {
		var role models.Role

		err := db.Where("role_name = ?", roleName).First(&role).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				newRole := models.Role{
					RoleName: roleName,
				}
				if err := db.Create(&newRole).Error; err != nil {
					return err
				}
				log.Println("Seeded role:", roleName)
			} else {
				return err
			}
		}
	}

	return nil
}

func SeedRolePermissionGroups(db *gorm.DB) error {

	// Roles
	var admin, member, guest models.Role
	if err := db.Where("role_name = ?", "admin").First(&admin).Error; err != nil {
		return err
	}
	if err := db.Where("role_name = ?", "member").First(&member).Error; err != nil {
		return err
	}
	if err := db.Where("role_name = ?", "guest").First(&guest).Error; err != nil {
		return err
	}

	// Permission Groups
	var userGroup, activityGroup models.PermissionGroup
	if err := db.Where("name = ?", UserManagementGroup).First(&userGroup).Error; err != nil {
		return err
	}
	if err := db.Where("name = ?", ActivityManagementGroup).First(&activityGroup).Error; err != nil {
		return err
	}

	// admin → user + activity
	if err := db.Model(&admin).
		Association("PermissionGroup").
		Replace(&userGroup, &activityGroup); err != nil {
		return err
	}
	log.Println("Admin assigned user & activity management")

	// member → user + activity
	if err := db.Model(&member).
		Association("PermissionGroup").
		Replace(&userGroup, &activityGroup); err != nil {
		return err
	}
	log.Println("Member assigned user & activity management")

	// guest → activity
	if err := db.Model(&guest).
		Association("PermissionGroup").
		Replace(&activityGroup); err != nil {
		return err
	}
	log.Println("Guest assigned activity management")

	return nil
}

func SeedRolePermissions(db *gorm.DB) error {

	// 1. ดึง Roles ที่ต้องการมาใช้งาน
	var admin, member, guest models.Role
	if err := db.Where("role_name = ?", "admin").First(&admin).Error; err != nil {
		return err
	}
	if err := db.Where("role_name = ?", "member").First(&member).Error; err != nil {
		return err
	}
	if err := db.Where("role_name = ?", "guest").First(&guest).Error; err != nil {
		return err
	}

	// 2. ดึง Permissions ที่จำเป็นทั้งหมดตามข้อกำหนด

	// **********************************************
	// ** ส่วนที่แก้ไข: เปลี่ยน Joins เป็น Raw SQL **
	// **********************************************

	// ดึงสิทธิ์ User Management ทั้งหมด
	var userPermissions []models.Permission
	if err := db.
		// ใช้ Raw SQL JOIN เพื่อเลี่ยงปัญหา Alias และ Case-Sensitivity
		Joins("LEFT JOIN permission_groups ON permissions.permission_group_id = permission_groups.id").
		Where("permission_groups.name = ?", UserManagementGroup).
		Find(&userPermissions).Error; err != nil {
		return err
	}

	// ดึงสิทธิ์ read_activity (ไม่ต้องแก้ไขส่วนนี้เพราะไม่ได้ใช้ Join)
	var readActivityPermission models.Permission
	if err := db.Where("permission_name = ?", "read_activity").First(&readActivityPermission).Error; err != nil {
		log.Printf("Warning: 'read_activity' permission not found. Error: %v", err)
		return err
	}

	// ดึงสิทธิ์ทั้งหมด (สำหรับ Admin)
	var allPermissions []models.Permission
	if err := db.Find(&allPermissions).Error; err != nil {
		return err
	}
	log.Printf("Found total %d permissions in the database.", len(allPermissions))

	// 3. กำหนดชุดสิทธิ์สำหรับแต่ละ Role

	// สิทธิ์สำหรับ Member: User Management ทั้งหมด + read_activity
	memberPermissions := append(userPermissions, readActivityPermission)

	// สิทธิ์สำหรับ Guest: read_activity เท่านั้น
	guestPermissions := []models.Permission{readActivityPermission}

	// 4. ผูก Permissions เข้ากับ Roles ด้วย GORM Association

	// 4.1. Admin ได้สิทธิ์ทั้งหมด
	if err := db.Model(&admin).
		Association("Permissions").
		Replace(allPermissions); err != nil {
		return err
	}
	log.Printf("Admin assigned %d total permissions.", len(allPermissions))

	// 4.2. Member ได้สิทธิ์ตามที่กำหนด
	if err := db.Model(&member).
		Association("Permissions").
		Replace(memberPermissions); err != nil {
		return err
	}
	log.Printf("Member assigned %d permissions (User Management + read_activity).", len(memberPermissions))

	// 4.3. Guest ได้แค่สิทธิ์ read_activity
	if err := db.Model(&guest).
		Association("Permissions").
		Replace(guestPermissions); err != nil {
		return err
	}
	log.Printf("Guest assigned %d permission (read_activity only).", len(guestPermissions))

	return nil
}

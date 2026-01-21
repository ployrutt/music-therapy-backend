package seeds

import (
	"log"
	"project-backend/models"

	"gorm.io/gorm"
)

func SeedPermissions(db *gorm.DB) error {

	// หา Permission Groups
	var userGroup models.PermissionGroup
	var activityGroup models.PermissionGroup

	if err := db.Where("name = ?", UserManagementGroup).
		First(&userGroup).Error; err != nil {
		return err
	}

	if err := db.Where("name = ?", ActivityManagementGroup).
		First(&activityGroup).Error; err != nil {
		return err
	}

	permissions := []models.Permission{
		// User Management
		{PermissionName: "create_user", PermissionGroupID: userGroup.ID},
		{PermissionName: "update_user", PermissionGroupID: userGroup.ID},
		{PermissionName: "delete_user", PermissionGroupID: userGroup.ID},
		{PermissionName: "view_user", PermissionGroupID: userGroup.ID},
		{PermissionName: "password_reset", PermissionGroupID: userGroup.ID},

		// Activity Management
		{PermissionName: "create_activity", PermissionGroupID: activityGroup.ID},
		{PermissionName: "update_activity", PermissionGroupID: activityGroup.ID},
		{PermissionName: "delete_activity", PermissionGroupID: activityGroup.ID},
		{PermissionName: "read_activity", PermissionGroupID: activityGroup.ID},
	}

	for _, p := range permissions {
		var exist models.Permission

		err := db.Where("permission_name = ?", p.PermissionName).
			First(&exist).Error

		if err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&p).Error; err != nil {
					return err
				}
				log.Println("Seeded Permission:", p.PermissionName)
			} else {
				return err
			}
		}
	}

	return nil
}

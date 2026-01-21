package seeds

import (
	"log"

	"project-backend/models"

	"gorm.io/gorm"
)

const (
	UserManagementGroup     = "user_management"
	ActivityManagementGroup = "activity_management"
)

func SeedPermissionGroups(db *gorm.DB) error {

	groups := []string{
		UserManagementGroup,
		ActivityManagementGroup,
	}

	for _, name := range groups {
		var group models.PermissionGroup

		err := db.Where("name = ?", name).First(&group).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&models.PermissionGroup{
					Name: name,
				}).Error; err != nil {
					return err
				}
				log.Println("Seeded Permission Group:", name)
			} else {
				return err
			}
		}
	}

	return nil
}

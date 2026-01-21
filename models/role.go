package models

import "time"

type Role struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement;index"`
	RoleName  string    `json:"role_name" gorm:"unique;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	User            []User            ` gorm:"foreignKey:RoleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	PermissionGroup []PermissionGroup `gorm:"many2many:role_permission_groups;"`
	Permissions     []Permission      `gorm:"many2many:role_permissions;"`
}

package models

import "time"

type Permission struct {
	ID             uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	PermissionName string `json:"permission_name" gorm:"unique;not null"`

	PermissionGroupID uint
	PermissionGroup   PermissionGroup

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

package models

import "time"

type PermissionGroup struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"permission_group_name" gorm:"unique;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	Permission []Permission `gorm:"foreignKey:PermissionGroupID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	FirstName   string         `json:"first_name" gorm:"column:firstname;not null"`
	LastName    string         `json:"last_name" gorm:"column:lastname;not null"`
	DateOfBirth time.Time      `json:"date_of_birth" gorm:"column:date_of_birth;type:date"`
	Email       string         `json:"email" gorm:"column:email;not null;unique"`
	Password    string         `json:"-" gorm:"column:password;not null"`
	PhoneNumber string         `json:"phone_number" gorm:"column:phone;not null"`
	Profile     string         `json:"profile" gorm:"column:profile"`
	RoleID      uint           `json:"role_id"`
	Role        *Role          `json:"-" gorm:"foreignKey:RoleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	Favorites   []UserFavorite    `json:"favorites" gorm:"foreignKey:UserID"`
	ReadHistory []UserReadHistory `json:"read_history" gorm:"foreignKey:UserID"`
}

type UserFavorite struct {
	UserID     uint      `json:"user_id" gorm:"column:user_id;primaryKey"`
	ActivityID uint      `json:"activity_id" gorm:"column:activity_id;primaryKey"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`

	User     User     `json:"-" gorm:"foreignKey:UserID;references:ID"`
	Activity Activity `json:"activity" gorm:"foreignKey:ActivityID;references:ID"`
}

type UserReadHistory struct {
	ID         uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID     uint      `json:"user_id" gorm:"index"`
	ActivityID uint      `json:"activity_id" gorm:"index"`
	ReadCount  int       `json:"read_count" gorm:"default:1"`
	UpdatedAt  time.Time `json:"last_read_at" gorm:"autoUpdateTime"`

	Activity Activity `json:"activity" gorm:"foreignKey:ActivityID"`
}

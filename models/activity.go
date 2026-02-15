package models

import "time"

type Activity struct {
	ID                 uint      `json:"activity_id" gorm:"primaryKey;autoIncrement"`
	Title              string    `json:"title" gorm:"type:text;not null"`
	CreatedAt          time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt          time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	CoverImage         string    `json:"cover_image" gorm:"type:text"`
	GoalDescription    string    `json:"goal_description" gorm:"type:text"`
	Equipment          string    `json:"equipment" gorm:"type:text"`
	Process            string    `json:"process" gorm:"type:text"`
	ObservableBehavior string    `json:"observable_behavior" gorm:"type:text"`
	Suggestion         string    `json:"suggestion" gorm:"type:text"`
	Song               string    `json:"song" gorm:"type:text"`
	SongImage          string    `json:"song_image" gorm:"type:text"`
	QR1                string    `json:"qr_1" gorm:"type:text"`
	QR2                string    `json:"qr_2" gorm:"type:text"`
	AdminID            uint      `json:"admin_id" gorm:"not null"`

	SubGoals      []ActivitySubGoal     `json:"selected_sub_goals" gorm:"many2many:activity_selected_sub_goals;"`
	SubCategories []ActivitySubCategory `json:"selected_sub_categories" gorm:"many2many:activity_selected_sub_categories;"`
}

type ActivityGoal struct {
	ID       uint              `json:"goal_id" gorm:"primaryKey"`
	GoalName string            `json:"goal_name" gorm:"type:text;not null"`
	SubGoals []ActivitySubGoal `json:"sub_goals" gorm:"foreignKey:GoalID"`
}

type ActivitySubGoal struct {
	ID          uint   `json:"sub_goal_id" gorm:"primaryKey"`
	GoalID      uint   `json:"goal_id"`
	SubGoalName string `json:"sub_goal_name" gorm:"type:text;not null"`
}

type ActivityMainCategory struct {
	ID           uint   `json:"category_id" gorm:"primaryKey"`
	CategoryName string `json:"category_name" gorm:"type:text;not null"`

	SubCategories []ActivitySubCategory `json:"sub_categories" gorm:"foreignKey:CategoryID"`
}

type ActivitySubCategory struct {
	ID uint `json:"sub_category_id" gorm:"primaryKey"`

	CategoryID      uint   `json:"category_id" gorm:"column:category_id"`
	SubCategoryName string `json:"sub_category_name" gorm:"type:text;not null"`
}

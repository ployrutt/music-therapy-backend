package controllers

import (
	"net/http"
	"strconv"

	"project-backend/models"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

func CreateActivity(db *gorm.DB) gin.HandlerFunc {

	return func(c *gin.Context) {
		var input struct {
			Title string `json:"title"`

			CoverImage string `json:"cover_image"`

			GoalDescription string `json:"goal_description"`

			Equipment string `json:"equipment"`

			Process string `json:"process"`

			ObservableBehavior string `json:"observable_behavior"`

			Suggestion string `json:"suggestion"`

			Song string `json:"song"`

			SongImage string `json:"song_image"`

			QR1 string `json:"qr_1"`

			QR2 string `json:"qr_2"`

			SubGoalIDs []uint `json:"sub_goal_ids"`

			SubCategoryIDs []uint `json:"sub_category_ids"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {

			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

			return

		}

		val, exists := c.Get("user_id")

		if !exists {

			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})

			return

		}

		userID := val.(uint)

		var selectedSubGoals []models.ActivitySubGoal

		db.Where("id IN ?", input.SubGoalIDs).Find(&selectedSubGoals)

		var selectedSubCats []models.ActivitySubCategory

		db.Where("id IN ?", input.SubCategoryIDs).Find(&selectedSubCats)

		activity := models.Activity{

			Title: input.Title,

			CoverImage: input.CoverImage,

			GoalDescription: input.GoalDescription,

			Equipment: input.Equipment,

			Process: input.Process,

			ObservableBehavior: input.ObservableBehavior,

			Suggestion: input.Suggestion,

			Song: input.Song,

			SongImage: input.SongImage,

			QR1: input.QR1,

			QR2: input.QR2,

			AdminID: userID,

			SubGoals: selectedSubGoals,

			SubCategories: selectedSubCats,
		}

		if err := db.Debug().Create(&activity).Error; err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

			return

		}
		c.JSON(http.StatusCreated, activity)
	}

}

func UpdateActivity(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var activity models.Activity
		if err := db.Preload("SubGoals").Preload("SubCategories").First(&activity, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
			return
		}
		var input struct {
			Title              string `json:"title"`
			GoalDescription    string `json:"goal_description"`
			SubGoalIDs         []uint `json:"sub_goal_ids"`
			SubCategoryIDs     []uint `json:"sub_category_ids"`
			CoverImage         string `json:"cover_image"`
			Equipment          string `json:"equipment"`
			Process            string `json:"process"`
			ObservableBehavior string `json:"observable_behavior"`
			Suggestion         string `json:"suggestion"`
			Song               string `json:"song"`
			SongImage          string `json:"song_image"`
			QR1                string `json:"qr_1"`
			QR2                string `json:"qr_2"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var newSubGoals []models.ActivitySubGoal
		db.Where("id IN ?", input.SubGoalIDs).Find(&newSubGoals)
		var newSubCats []models.ActivitySubCategory
		db.Where("id IN ?", input.SubCategoryIDs).Find(&newSubCats)
		err := db.Transaction(func(tx *gorm.DB) error {
			updates := map[string]interface{}{
				"title":               input.Title,
				"goal_description":    input.GoalDescription,
				"cover_image":         input.CoverImage,
				"equipment":           input.Equipment,
				"process":             input.Process,
				"observable_behavior": input.ObservableBehavior,
				"suggestion":          input.Suggestion,
				"song":                input.Song,
				"song_image":          input.SongImage,
				"qr1":                 input.QR1,
				"qr2":                 input.QR2,
			}
			if err := tx.Model(&activity).Updates(updates).Error; err != nil {
				return err
			}
			if err := tx.Model(&activity).Association("SubGoals").Replace(newSubGoals); err != nil {
				return err
			}
			if err := tx.Model(&activity).Association("SubCategories").Replace(newSubCats); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed: " + err.Error()})
			return
		}
		db.Preload("SubGoals").Preload("SubCategories").First(&activity, id)
		c.JSON(http.StatusOK, activity)
	}
}

func DeleteActivity(db *gorm.DB) gin.HandlerFunc {

	return func(c *gin.Context) {

		id := c.Param("id")

		var activity models.Activity

		if err := db.Preload("SubGoals").Preload("SubCategories").First(&activity, id).Error; err != nil {

			c.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})

			return

		}

		err := db.Transaction(func(tx *gorm.DB) error {

			if err := tx.Model(&activity).Association("SubGoals").Clear(); err != nil {

				return err

			}

			if err := tx.Model(&activity).Association("SubCategories").Clear(); err != nil {

				return err

			}

			return tx.Delete(&activity).Error

		})

		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

			return

		}

		c.JSON(http.StatusOK, gin.H{"message": "Activity deleted successfully"})

	}

}

func GetActivityByID(db *gorm.DB) gin.HandlerFunc {

	return func(c *gin.Context) {

		id := c.Param("id")

		var activity models.Activity

		if err := db.
			Preload("SubGoals").
			Preload("SubCategories").
			First(&activity, id).Error; err != nil {

			c.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})

			return

		}

		c.JSON(http.StatusOK, activity)

	}

}

func ListActivities(db *gorm.DB) gin.HandlerFunc {

	return func(c *gin.Context) {

		var activities []models.Activity

		if err := db.
			Preload("SubGoals").
			Preload("SubCategories").
			Order("id DESC").
			Find(&activities).Error; err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

			return

		}

		c.JSON(http.StatusOK, activities)

	}

}

func GetActivityMasterGoals(db *gorm.DB) gin.HandlerFunc {

	return func(c *gin.Context) {

		var goals []models.ActivityGoal

		if err := db.Preload("SubGoals").Order("id ASC").Find(&goals).Error; err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

			return

		}

		c.JSON(http.StatusOK, goals)

	}

}

func GetActivityMasterCategories(db *gorm.DB) gin.HandlerFunc {

	return func(c *gin.Context) {

		var categories []models.ActivityMainCategory

		if err := db.Preload("SubCategories").Order("id ASC").Find(&categories).Error; err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch categories: " + err.Error()})

			return

		}

		c.JSON(http.StatusOK, categories)

	}

}

func ToggleFavorite(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		val, _ := c.Get("user_id")
		userID := val.(uint)

		idStr := c.Param("id")
		activityID, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
			return
		}

		var fav models.UserFavorite

		result := db.Where("user_id = ? AND activity_id = ?", userID, uint(activityID)).First(&fav)

		if result.Error == nil {

			db.Delete(&fav)
			c.JSON(http.StatusOK, gin.H{"status": "unfavorited", "message": "ลบออกจากรายการโปรดแล้ว"})
		} else {

			newFav := models.UserFavorite{
				UserID:     userID,
				ActivityID: uint(activityID),
			}
			if err := db.Create(&newFav).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, gin.H{"status": "favorited", "message": "เพิ่มในรายการโปรดแล้ว"})
		}
	}
}

func ListFavorites(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, _ := c.Get("user_id")
		userID := val.(uint)

		var favorites []models.UserFavorite

		err := db.Preload("Activity", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "title")
		}).Where("user_id = ?", userID).Find(&favorites).Error

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, favorites)
	}
}

func RecordReadHistory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, _ := c.Get("user_id")
		userID := val.(uint)

		idStr := c.Param("id")
		activityID, _ := strconv.ParseUint(idStr, 10, 32)

		var history models.UserReadHistory
		result := db.Where("user_id = ? AND activity_id = ?", userID, uint(activityID)).First(&history)

		if result.Error != nil {

			newHistory := models.UserReadHistory{
				UserID:     userID,
				ActivityID: uint(activityID),
				ReadCount:  1,
			}
			db.Create(&newHistory)
		} else {

			db.Model(&history).Update("read_count", gorm.Expr("read_count + ?", 1))
		}
		c.JSON(http.StatusOK, gin.H{"message": "บันทึกประวัติการอ่านเรียบร้อย"})
	}
}
func ListReadHistory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, _ := c.Get("user_id")
		userID := val.(uint)

		var history []models.UserReadHistory

		err := db.Preload("Activity", func(db *gorm.DB) *gorm.DB {

			return db.Select("id", "title")
		}).Where("user_id = ?", userID).
			Order("updated_at DESC").
			Find(&history).Error

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, history)
	}
}

func SearchAndFilterActivities(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var activities []models.Activity
		query := db.Model(&models.Activity{}).Preload("SubGoals").Preload("SubCategories")

		if title := c.Query("title"); title != "" {
			query = query.Where("title LIKE ?", "%"+title+"%")
		}

		if subGoalID := c.Query("sub_goal_id"); subGoalID != "" {
			query = query.Joins("JOIN activity_sub_goal_assoc ON activity_sub_goal_assoc.activity_id = activities.id").
				Where("activity_sub_goal_assoc.activity_sub_goal_id = ?", subGoalID)
		}

		if subCatID := c.Query("sub_category_id"); subCatID != "" {
			query = query.Joins("JOIN activity_sub_category_assoc ON activity_sub_category_assoc.activity_id = activities.id").
				Where("activity_sub_category_assoc.activity_sub_category_id = ?", subCatID)
		}

		if err := query.Distinct().Find(&activities).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, activities)
	}
}

func GetActivityStats(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var favCount int64
		var readTotal int64

		db.Model(&models.UserFavorite{}).Where("activity_id = ?", id).Count(&favCount)

		db.Model(&models.UserReadHistory{}).Where("activity_id = ?", id).Select("COALESCE(sum(read_count), 0)").Row().Scan(&readTotal)

		c.JSON(http.StatusOK, gin.H{
			"activity_id":    id,
			"favorite_count": favCount,
			"total_reads":    readTotal,
		})
	}
}

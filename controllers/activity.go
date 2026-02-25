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

// func SearchAndFilterActivities(db *gorm.DB) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var activities []models.Activity
// 		query := db.Model(&models.Activity{}).Preload("SubGoals").Preload("SubCategories")

// 		if title := c.Query("title"); title != "" {
// 			query = query.Where("title LIKE ?", "%"+title+"%")
// 		}

// 		if subGoalID := c.Query("sub_goal_id"); subGoalID != "" {
// 			query = query.Joins("JOIN activity_sub_goal_assoc ON activity_sub_goal_assoc.activity_id = activities.id").
// 				Where("activity_sub_goal_assoc.activity_sub_goal_id = ?", subGoalID)
// 		}

// 		if subCatID := c.Query("sub_category_id"); subCatID != "" {
// 			query = query.Joins("JOIN activity_sub_category_assoc ON activity_sub_category_assoc.activity_id = activities.id").
// 				Where("activity_sub_category_assoc.activity_sub_category_id = ?", subCatID)
// 		}

// 		if err := query.Distinct().Find(&activities).Error; err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}
// 		c.JSON(http.StatusOK, activities)
// 	}
// }

func SearchAndFilterActivities(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var activities []models.Activity
		// เริ่มต้น Query และ Preload ข้อมูลที่เกี่ยวข้องมาแสดงผลด้วย
		query := db.Model(&models.Activity{}).Preload("SubGoals").Preload("SubCategories")

		// 1. ค้นหาจากชื่อ (Title)
		if title := c.Query("title"); title != "" {
			query = query.Where("activities.title LIKE ?", "%"+title+"%")
		}

		// 2. ค้นหาจาก เป้าหมายย่อย (Sub Goal)
		if subGoalID := c.Query("sub_goal_id"); subGoalID != "" {
			query = query.Joins("JOIN activity_selected_sub_goals ON activity_selected_sub_goals.activity_id = activities.id").
				Where("activity_selected_sub_goals.activity_sub_goal_id = ?", subGoalID)
		}

		// 3. ค้นหาจาก เป้าหมายหลัก (Master Goal) [เพิ่มใหม่]
		// ต้อง JOIN 2 ต่อ: Activity -> SubGoalAssoc -> SubGoal
		if goalID := c.Query("goal_id"); goalID != "" {
			query = query.Joins("JOIN activity_selected_sub_goals ON activity_selected_sub_goals.activity_id = activities.id").
				Joins("JOIN activity_sub_goals ON activity_sub_goals.id = activity_selected_sub_goals.activity_sub_goal_id").
				Where("activity_sub_goals.goal_id = ?", goalID)
		}

		// 4. ค้นหาจาก หมวดหมู่ย่อย (Sub Category)
		if subCatID := c.Query("sub_category_id"); subCatID != "" {
			query = query.Joins("JOIN activity_selected_sub_categories ON activity_selected_sub_categories.activity_id = activities.id").
				Where("activity_selected_sub_categories.activity_sub_category_id = ?", subCatID)
		}

		// 5. ค้นหาจาก หมวดหมู่หลัก (Master Category) [เพิ่มใหม่]
		// ต้อง JOIN 2 ต่อ: Activity -> SubCategoryAssoc -> SubCategory
		if catID := c.Query("category_id"); catID != "" {
			query = query.Joins("JOIN activity_selected_sub_categories ON activity_selected_sub_categories.activity_id = activities.id").
				Joins("JOIN activity_sub_categories ON activity_sub_categories.id = activity_selected_sub_categories.activity_sub_category_id").
				Where("activity_sub_categories.category_id = ?", catID)
		}

		// ใช้ .Distinct() เพื่อป้องกันข้อมูลซ้ำกรณีที่ 1 กิจกรรมมีหลาย Sub-goal ใน Master เดียวกัน
		if err := query.Distinct("activities.*").Order("activities.id DESC").Find(&activities).Error; err != nil {
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

func GetAdminDashboard(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// รับค่าช่วงเวลาจาก Query Params (เช่น ?start=2024-01-01&end=2024-01-31)
		startDate := c.Query("start")
		endDate := c.Query("end")

		// สร้าง Base Query สำหรับการกรองช่วงเวลา (ถ้ามี)
		filterRead := db.Model(&models.UserReadHistory{})
		if startDate != "" && endDate != "" {
			filterRead = filterRead.Where("updated_at BETWEEN ? AND ?", startDate, endDate)
		}

		// 1. กิจกรรมที่มีคนอ่านมากที่สุด (Top Read Activities)
		var topReadActivities []struct {
			ActivityID uint   `json:"activity_id"`
			Title      string `json:"title"`
			TotalRead  int    `json:"total_read"`
		}
		filterRead.Select("activity_id, activities.title, sum(read_count) as total_read").
			Joins("JOIN activities ON activities.id = user_read_histories.activity_id").
			Group("activity_id, activities.title").
			Order("total_read DESC").Limit(10).Scan(&topReadActivities)

		// 2. กิจกรรมที่มีคนกด Favorite มากที่สุด (Top Favorited Activities)
		var topFavActivities []struct {
			ActivityID uint   `json:"activity_id"`
			Title      string `json:"title"`
			FavCount   int    `json:"fav_count"`
		}
		db.Model(&models.UserFavorite{}).
			Select("activity_id, activities.title, count(*) as fav_count").
			Joins("JOIN activities ON activities.id = user_favorites.activity_id").
			Group("activity_id, activities.title").
			Order("fav_count DESC").Limit(10).Scan(&topFavActivities)

		// 3. หมวดหมู่ย่อยและ Goals ที่มีคน Favorite มากที่สุด
		var topFavCategories []struct {
			SubCategoryName string `json:"name"`
			Count           int    `json:"count"`
		}
		db.Table("user_favorites").
			Select("activity_sub_categories.sub_category_name, count(*) as count").
			Joins("JOIN activity_selected_sub_categories ON activity_selected_sub_categories.activity_id = user_favorites.activity_id").
			Joins("JOIN activity_sub_categories ON activity_sub_categories.id = activity_selected_sub_categories.activity_sub_category_id").
			Group("activity_sub_categories.sub_category_name").Order("count DESC").Limit(5).Scan(&topFavCategories)

		// 4. หมวดหมู่และ Goals ที่มีคนอ่านมากที่สุด (ตามช่วงเวลา)
		var topReadGoals []struct {
			SubGoalName string `json:"name"`
			TotalRead   int    `json:"total_read"`
		}
		filterRead.Table("user_read_histories").
			Select("activity_sub_goals.sub_goal_name, sum(user_read_histories.read_count) as total_read").
			Joins("JOIN activity_selected_sub_goals ON activity_selected_sub_goals.activity_id = user_read_histories.activity_id").
			Joins("JOIN activity_sub_goals ON activity_sub_goals.id = activity_selected_sub_goals.activity_sub_goal_id").
			Group("activity_sub_goals.sub_goal_name").Order("total_read DESC").Limit(5).Scan(&topReadGoals)

		c.JSON(http.StatusOK, gin.H{
			"top_read_activities": topReadActivities,
			"top_fav_activities":  topFavActivities,
			"top_fav_categories":  topFavCategories,
			"top_read_goals":      topReadGoals,
		})
	}
}

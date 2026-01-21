package seeds

import (
	"log"
	"project-backend/models"

	"gorm.io/gorm"
)

const (
	CreateActivity = "create_activity"
	UpdateActivity = "update_activity"
	DeleteActivity = "delete_activity"
	ViewActivity   = "view_activity"
)

func SeedActivityPermissions(db *gorm.DB) error {
	var group models.PermissionGroup
	if err := db.Where("name = ?", ActivityManagementGroup).First(&group).Error; err != nil {
		return err
	}

	perms := []string{CreateActivity, UpdateActivity, DeleteActivity, ViewActivity}

	for _, name := range perms {
		var p models.Permission
		err := db.Where("permission_name = ?", name).First(&p).Error

		if err == gorm.ErrRecordNotFound {
			if err := db.Create(&models.Permission{
				PermissionName:    name,
				PermissionGroupID: group.ID,
			}).Error; err != nil {
				return err
			}
			log.Println("Seeded Permission:", name)
		}
	}
	return nil
}

func SeedActivityGoals(db *gorm.DB) error {
	goals := []models.ActivityGoal{
		{
			ID:       1,
			GoalName: "พัฒนาการทางสติปัญญา",
			SubGoals: []models.ActivitySubGoal{
				{ID: 1, SubGoalName: "จดจ่อสิ่งใดสิ่งหนึ่งเป็นระยะเวลาที่เหมาะสม"},
				{ID: 2, SubGoalName: "จดจำทำนองหรือรูปแบบจังหวะ"},
				{ID: 3, SubGoalName: "ระบุที่มาของเสียง"},
				{ID: 4, SubGoalName: "แยกแยะเสียงสูง-ต่ำ/ช้า-เร็ว"},
				{ID: 5, SubGoalName: "จดจำและเลียนแบบท่าทาง"},
				{ID: 6, SubGoalName: "แสดงความคิดสร้างสรรค์"},
			},
		},
		{
			ID:       2,
			GoalName: "พัฒนาการทางร่างกาย",
			SubGoals: []models.ActivitySubGoal{
				{ID: 7, SubGoalName: "การเคลื่อนไหวตามจังหวะอย่างสม่ำเสมอ"},
				{ID: 8, SubGoalName: "การควบคุมร่างกายในการใช้เครื่องดนตรี"},
				{ID: 9, SubGoalName: "รับ-ส่งเครื่องดนตรีและอุปกรณ์"},
			},
		},
		{
			ID:       3,
			GoalName: "พัฒนาการทางอารมณ์",
			SubGoals: []models.ActivitySubGoal{
				{ID: 10, SubGoalName: "อดทนต่อการรอคอย"},
				{ID: 11, SubGoalName: "อดทนต่อปัญหาหรือสิ่งที่ตนเองไม่ชอบ"},
				{ID: 12, SubGoalName: "แสดงออกทางอารมณ์อย่างเหมาะสม"},
				{ID: 13, SubGoalName: "มีความมั่นใจและกล้าแสดงออก"},
			},
		},
		{
			ID:       4,
			GoalName: "พัฒนาการทางสังคม",
			SubGoals: []models.ActivitySubGoal{
				{ID: 14, SubGoalName: "การสบตาหรือการมองผู้อื่น"},
				{ID: 15, SubGoalName: "การทำตามสัญญาณหรือข้อเสนอแนะของผู้อื่น"},
				{ID: 16, SubGoalName: "การแสดงความคิดเห็นของตนเอง"},
				{ID: 17, SubGoalName: "การเป็นผู้นำ"},
				{ID: 18, SubGoalName: "การสลับกันเล่น"},
				{ID: 19, SubGoalName: "แบ่งเครื่องดนตรีหรืออุปกรณ์กับผู้อื่น"},
				{ID: 20, SubGoalName: "แสดงความสนใจร่วมกับผู้อื่น"},
				{ID: 21, SubGoalName: "ตระหนักรู้ว่ามีผู้อื่นอยู่"},
			},
		},
	}

	for _, goal := range goals {
		// ใช้ Create เพื่อให้ FullSaveAssociations ทำงานได้อย่างถูกต้อง
		if err := db.Create(&goal).Error; err != nil {
			log.Printf("Skip seeding Goal %d: %v", goal.ID, err)
		}
	}
	return nil
}

func SeedMainCategories(db *gorm.DB) error {
	categories := []models.ActivityMainCategory{
		{
			ID:           1,
			CategoryName: "การร้องเพลง",
			SubCategories: []models.ActivitySubCategory{
				{ID: 1, SubCategoryName: "การร้องเพลงเดี่ยว"},
				{ID: 2, SubCategoryName: "การร้องเพลงกลุ่มย่อย"},
				{ID: 3, SubCategoryName: "การร้องเพลงร่วมกับทุกคน"},
				{ID: 4, SubCategoryName: "การสลับกันออกเสียงอย่างอิสระ"},
				{ID: 5, SubCategoryName: "ออกเสียงเลียนแบบสิ่งรอบตัว"},
			},
		},
		{
			ID:           2,
			CategoryName: "การฟัง",
			SubCategories: []models.ActivitySubCategory{
				{ID: 6, SubCategoryName: "แยกแยะสัญญาณแทนสิ่งรอบตัว"},
				{ID: 7, SubCategoryName: "แยะแยะเสียงดนตรีที่เหมือน-ต่าง"},
				{ID: 8, SubCategoryName: "จดจำสัญญาณเสียง"},
				{ID: 9, SubCategoryName: "แยะแยะเสียงดนตรีดัง-เบา/ช้า-เร็ว/สูง-ต่ำ"},
			},
		},
		{
			ID:           3,
			CategoryName: "การเคลื่อนไหว",
			SubCategories: []models.ActivitySubCategory{
				{ID: 10, SubCategoryName: "เคลื่อนไหวตามสัญญาณ"},
				{ID: 11, SubCategoryName: "เคลื่อนไหวเลียนแบบสิ่งรอบตัว"},
				{ID: 12, SubCategoryName: "เคลื่อนไหวอย่างอิสระ"},
			},
		},
		{
			ID:           4,
			CategoryName: "การบรรเลงเครื่องดนตรี",
			SubCategories: []models.ActivitySubCategory{
				{ID: 13, SubCategoryName: "บรรเลงอิสระกับกลุ่มย่อย"},
				{ID: 14, SubCategoryName: "บรรเลงอย่างอิสระร่วมกับทุกคน"},
				{ID: 15, SubCategoryName: "บรรเลงดนตรีร่วมกับผู้อื่นตามสัญญาณ"},
				{ID: 16, SubCategoryName: "การใช้ร่างกายสร้างเสียง(Body Percussion)"},
				{ID: 17, SubCategoryName: "บรรเลงดนตรีอย่างอิสระ(ทีละคน)"},
				{ID: 18, SubCategoryName: "บรรเลงเครื่องดนตรีประกอบการเคลื่อนไหว"},
			},
		},
	}

	for _, cat := range categories {
		// ใช้ Session และระบุ FullSaveAssociations: true
		// เพื่อบังคับให้ GORM ไล่เก็บข้อมูลใน []ActivitySubCategory ลงตารางด้วย
		if err := db.Session(&gorm.Session{FullSaveAssociations: true}).Create(&cat).Error; err != nil {
			log.Printf("❌ Error seeding Category %d: %v", cat.ID, err)
		} else {
			log.Printf("✅ Seeded Category %d with %d sub-categories", cat.ID, len(cat.SubCategories))
		}
	}
	return nil
}

// func SeedActivityGoals(db *gorm.DB) error {
// 	goals := []models.ActivityGoal{
// 		{
// 			ID:       1,
// 			GoalName: "พัฒนาการทางสติปัญญา",
// 			SubGoals: []models.ActivitySubGoal{
// 				{ID: 1, SubGoalName: "จดจ่อสิ่งใดสิ่งหนึ่งเป็นระยะเวลาที่เหมาะสม"},
// 				{ID: 2, SubGoalName: "จดจำทำนองหรือรูปแบบจังหวะ"},
// 				{ID: 3, SubGoalName: "ระบุที่มาของเสียง"},
// 				{ID: 4, SubGoalName: "แยกแยะเสียงสูง-ต่ำ/ช้า-เร็ว"},
// 				{ID: 5, SubGoalName: "จดจำและเลียนแบบท่าทาง"},
// 				{ID: 6, SubGoalName: "แสดงความคิดสร้างสรรค์"},
// 			},
// 		},
// 		{
// 			ID:       2,
// 			GoalName: "พัฒนาการทางร่างกาย",
// 			SubGoals: []models.ActivitySubGoal{
// 				{ID: 7, SubGoalName: "การเคลื่อนไหวตามจังหวะอย่างสม่ำเสมอ"},
// 				{ID: 8, SubGoalName: "การควบคุมร่างกายในการใช้เครื่องดนตรี"},
// 				{ID: 9, SubGoalName: "รับ-ส่งเครื่องดนตรีและอุปกรณ์"},
// 			},
// 		},
// 		{
// 			ID:       3,
// 			GoalName: "พัฒนาการทางอารมณ์",
// 			SubGoals: []models.ActivitySubGoal{
// 				{ID: 10, SubGoalName: "อดทนต่อการรอคอย"},
// 				{ID: 11, SubGoalName: "อดทนต่อปัญหาหรือสิ่งที่ตนเองไม่ชอบ"},
// 				{ID: 12, SubGoalName: "แสดงออกทางอารมณ์อย่างเหมาะสม"},
// 				{ID: 13, SubGoalName: "มีความมั่นใจและกล้าแสดงออก"},
// 			},
// 		},
// 		{
// 			ID:       4,
// 			GoalName: "พัฒนาการทางสังคม",
// 			SubGoals: []models.ActivitySubGoal{
// 				{ID: 14, SubGoalName: "การสบตาหรือการมองผู้อื่น"},
// 				{ID: 15, SubGoalName: "การทำตามสัญญาณหรือข้อเสนอแนะของผู้อื่น"},
// 				{ID: 16, SubGoalName: "การแสดงความคิดเห็นของตนเอง"},
// 				{ID: 17, SubGoalName: "การเป็นผู้นำ"},
// 				{ID: 18, SubGoalName: "การสลับกันเล่น"},
// 				{ID: 19, SubGoalName: "แบ่งเครื่องดนตรีหรืออุปกรณ์กับผู้อื่น"},
// 				{ID: 20, SubGoalName: "แสดงความสนใจร่วมกับผู้อื่น"},
// 				{ID: 21, SubGoalName: "ตระหนักรู้ว่ามีผู้อื่นอยู่"},
// 			},
// 		},
// 	}

// 	for _, goal := range goals {

// 		if err := db.Session(&gorm.Session{FullSaveAssociations: true}).FirstOrCreate(&models.ActivityGoal{}, goal).Error; err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }
// func SeedMainCategories(db *gorm.DB) error {

// 	categories := []models.ActivityMainCategory{
// 		{
// 			ID:           1,
// 			CategoryName: "การร้องเพลง",
// 			SubCategories: []models.ActivitySubCategory{
// 				{ID: 1, SubCategoryName: "การร้องเพลงเดี่ยว"},
// 				{ID: 2, SubCategoryName: "การร้องเพลงกลุ่มย่อย"},
// 				{ID: 3, SubCategoryName: "การร้องเพลงร่วมกับทุกคน"},
// 				{ID: 4, SubCategoryName: "การสลับกันออกเสียงอย่างอิสระ"},
// 				{ID: 5, SubCategoryName: "ออกเสียงเลียนแบบสิ่งรอบตัว"},
// 			},
// 		},
// 		{
// 			ID:           2,
// 			CategoryName: "การฟัง",
// 			SubCategories: []models.ActivitySubCategory{
// 				{ID: 6, SubCategoryName: "แยกแยะสัญญาณแทนสิ่งรอบตัว"},
// 				{ID: 7, SubCategoryName: "แยะแยะเสียงดนตรีที่เหมือน-ต่าง"},
// 				{ID: 8, SubCategoryName: "จดจำสัญญาณเสียง"},
// 				{ID: 9, SubCategoryName: "แยะแยะเสียงดนตรีดัง-เบา/ช้า-เร็ว/สูง-ต่ำ"},
// 			},
// 		},
// 		{
// 			ID:           3,
// 			CategoryName: "การเคลื่อนไหว",
// 			SubCategories: []models.ActivitySubCategory{
// 				{ID: 10, SubCategoryName: "เคลื่อนไหวตามสัญญาณ"},
// 				{ID: 11, SubCategoryName: "เคลื่อนไหวเลียนแบบสิ่งรอบตัว"},
// 				{ID: 12, SubCategoryName: "เคลื่อนไหวอย่างอิสระ"},
// 			},
// 		},
// 		{
// 			ID:           4,
// 			CategoryName: "การบรรเลงเครื่องดนตรี",
// 			SubCategories: []models.ActivitySubCategory{
// 				{ID: 13, SubCategoryName: "บรรเลงอิสระกับกลุ่มย่อย"},
// 				{ID: 14, SubCategoryName: "บรรเลงอย่างอิสระร่วมกับทุกคน"},
// 				{ID: 15, SubCategoryName: "บรรเลงดนตรีร่วมกับผู้อื่นตามสัญญาณ"},
// 				{ID: 16, SubCategoryName: "การใช้ร่างกายสร้างเสียง(Body Percussion)"},
// 				{ID: 17, SubCategoryName: "บรรเลงดนตรีอย่างอิสระ(ทีละคน)"},
// 				{ID: 18, SubCategoryName: "บรรเลงเครื่องดนตรีประกอบการเคลื่อนไหว"},
// 			},
// 		},
// 	}

// 	for _, cat := range categories {
// 		if err := db.Session(&gorm.Session{FullSaveAssociations: true}).FirstOrCreate(&models.ActivityMainCategory{}, cat).Error; err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

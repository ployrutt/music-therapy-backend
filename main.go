// package main

// import (
// 	"log"
// 	"time"

// 	"project-backend/config"
// 	"project-backend/db"
// 	"project-backend/models"
// 	"project-backend/router"
// 	"project-backend/seeds"

// 	"github.com/joho/godotenv"
// 	"golang.org/x/crypto/bcrypt"
// 	"gorm.io/gorm"
// )

// func hashPassword(password string) (string, error) {
// 	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return "", err
// 	}
// 	return string(bytes), nil
// }

// func main() {
// 	// โหลด .env file (ถ้ามี) - ใน Docker จะใช้ env vars จาก Dokploy แทน
// 	if err := godotenv.Load(); err != nil {
// 		log.Println("No .env file found, using system environment variables")
// 	}

// 	cfg := config.GetDBConfig()

// 	// 1. เชื่อมต่อ Database
// 	gormDB, err := db.InitDB(
// 		cfg,
// 		&models.User{},
// 		&models.Role{},
// 		&models.PermissionGroup{},
// 		&models.Permission{},
// 		&models.Activity{},
// 		&models.ActivityGoal{},
// 		&models.ActivitySubGoal{},
// 		&models.ActivityMainCategory{},
// 		&models.ActivitySubCategory{},

// 		&models.UserFavorite{},
// 		&models.UserReadHistory{},
// 	)

// 	if err != nil {
// 		log.Fatalf("Application startup failed: %v", err)
// 	}

// 	// 2. ล้างตารางเดิมทิ้ง (สำคัญมากสำหรับ PostgreSQL เพื่อไม่ให้ติด FK Error)
// 	log.Println("Resetting database tables...")
// 	gormDB.Exec(`DROP TABLE IF EXISTS
// 		role_permissions, role_permission_groups,
// 		activity_selected_sub_goals, activity_selected_sub_categories,
// 		activity_sub_goals, activity_sub_categories,
// 		activity_goals, activity_main_categories,
// 		activities, users, permissions, permission_groups, roles CASCADE`)

// 	// 3. สร้างตารางใหม่ทั้งหมด
// 	log.Println("Migrating database structure...")

// 	err = gormDB.AutoMigrate(
// 		&models.Role{},
// 		&models.PermissionGroup{},
// 		&models.Permission{},
// 		&models.User{},
// 		&models.Activity{},
// 		&models.ActivityGoal{},
// 		&models.ActivitySubGoal{},
// 		&models.ActivityMainCategory{},
// 		&models.ActivitySubCategory{},
// 	)
// 	if err != nil {
// 		log.Fatalf("Migration failed: %v", err)
// 	}

// 	// 4. รันข้อมูลเริ่มต้น
// 	runDatabaseSeeds(gormDB)

// 	log.Println("Database connection and migration successful.")

// 	// 5. เริ่มต้น Server
// 	r := router.SetupRouter(gormDB)
// 	port := ":8080"
// 	log.Printf("Starting HTTP server on port %s", port)

// 	sqlDB, _ := gormDB.DB()
// 	defer sqlDB.Close()

// 	if err := r.Run(port); err != nil {
// 		log.Fatalf("Server failed to start: %v", err)
// 	}
// }
// func runDatabaseSeeds(gormDB *gorm.DB) {
// 	log.Println("Starting database seeding...")

// 	seeds.SeedRoles(gormDB)
// 	seeds.SeedPermissionGroups(gormDB)
// 	seeds.SeedPermissions(gormDB)
// 	seeds.SeedRolePermissionGroups(gormDB)
// 	seeds.SeedRolePermissions(gormDB)

// 	if err := seeds.SeedActivityGoals(gormDB); err != nil {
// 		log.Printf("Error seeding ActivityGoals: %v", err)
// 	}
// 	if err := seeds.SeedMainCategories(gormDB); err != nil {
// 		log.Printf("Error seeding MainCategories: %v", err)
// 	}

// 	var adminCount int64

// 	gormDB.Model(&models.User{}).Where("email = ?", "Test@example.com").Count(&adminCount)

// 	if adminCount == 0 {

// 		hashedPass, _ := hashPassword("12052546")
// 		adminUser := models.User{
// 			FirstName:   "test",
// 			LastName:    "admin",
// 			DateOfBirth: time.Date(1999, time.March, 1, 0, 0, 0, 0, time.UTC),
// 			Email:       "TestAdmin@example.com",
// 			Password:    hashedPass,
// 			PhoneNumber: "0123456789",
// 			Profile:     "image",
// 			RoleID:      1,
// 		}
// 		if err := gormDB.Create(&adminUser).Error; err != nil {
// 			log.Printf("Could not create initial Admin: %v", err)
// 		} else {
// 			log.Println("Admin User created successfully.")
// 		}
// 	} else {
// 		log.Println("Admin User already exists, skipping...")
// 	}

//		log.Println("Seeding process completed.")
//	}
package main

import (
	"log"
	"os"
	"time"

	"project-backend/config"
	"project-backend/db"
	"project-backend/models"
	"project-backend/router"
	"project-backend/seeds"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func main() {
	// 1. โหลด .env เฉพาะตอนรัน Local
	// บน Dokploy จะหาไฟล์ไม่เจอ ซึ่งไม่เป็นไร ระบบจะไปอ่านจาก Environment Tab แทน
	if err := godotenv.Load(); err != nil {
		log.Println("Dokploy Environment detected: Using system variables")
	}

	// 2. ดึงค่า Config (จะใช้ DB_HOST=music-therapy-database-dlpqho อัตโนมัติ)
	cfg := config.GetDBConfig()

	// 3. เชื่อมต่อ Database และทำ AutoMigrate
	// เราจะไม่ใช้ DROP TABLE ในเวอร์ชันนี้เพื่อป้องกันข้อมูลบน Dokploy หาย
	gormDB, err := db.InitDB(
		cfg,
		&models.Role{}, // 1. มาก่อนเพื่อน
		&models.PermissionGroup{},
		&models.Permission{},
		&models.User{},     // 2. มาก่อน Favorite/ReadHistory
		&models.Activity{}, // 3. มาก่อน Favorite/ReadHistory
		&models.ActivityGoal{},
		&models.ActivitySubGoal{},
		&models.ActivityMainCategory{},
		&models.ActivitySubCategory{},

		// 4. สองตัวนี้อยู่ล่างสุดเสมอ เพราะต้องรอ User และ Activity
		&models.UserFavorite{},
		&models.UserReadHistory{},
	)

	if err != nil {
		log.Fatalf("Application startup failed: %v", err)
	}

	// 4. รันข้อมูลเริ่มต้น (Seeding)
	// ในฟังก์ชันนี้เราจะเพิ่ม Logic เช็คก่อนสร้าง เพื่อไม่ให้ข้อมูลซ้ำ
	runDatabaseSeeds(gormDB)

	log.Println("Database connection and migration successful.")

	// 5. ตั้งค่า Port (Dokploy มักจะส่ง PORT มาให้ หรือใช้ 8080 เป็นค่ามาตรฐาน)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 6. เริ่มต้น Server
	r := router.SetupRouter(gormDB)
	log.Printf("Starting HTTP server on port %s in %s mode", port, os.Getenv("GIN_MODE"))

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func runDatabaseSeeds(gormDB *gorm.DB) {
	log.Println("Starting database seeding...")

	// รัน Seed ข้อมูลพื้นฐาน
	seeds.SeedRoles(gormDB)
	seeds.SeedPermissionGroups(gormDB)
	seeds.SeedPermissions(gormDB)
	seeds.SeedRolePermissionGroups(gormDB)
	seeds.SeedRolePermissions(gormDB)

	if err := seeds.SeedActivityGoals(gormDB); err != nil {
		log.Printf("Error seeding ActivityGoals: %v", err)
	}
	if err := seeds.SeedMainCategories(gormDB); err != nil {
		log.Printf("Error seeding MainCategories: %v", err)
	}

	// --- ตรวจสอบและสร้าง Admin User ---
	var adminCount int64
	targetEmail := "TestAdmin@example.com" // ใช้อีเมลนี้เป็นหลัก

	// แก้ไข: เช็คจาก targetEmail ที่เราจะใช้สร้างจริงๆ
	gormDB.Model(&models.User{}).Where("email = ?", targetEmail).Count(&adminCount)

	if adminCount == 0 {
		hashedPass, _ := hashPassword("12052546")
		adminUser := models.User{
			FirstName:   "test",
			LastName:    "admin",
			DateOfBirth: time.Date(1999, time.March, 1, 0, 0, 0, 0, time.UTC),
			Email:       targetEmail, // ใช้ค่าเดียวกัน
			Password:    hashedPass,
			PhoneNumber: "0123456789",
			Profile:     "image",
			RoleID:      1,
		}
		if err := gormDB.Create(&adminUser).Error; err != nil {
			log.Printf("Could not create initial Admin: %v", err)
		} else {
			log.Println("Initial Admin User created successfully.")
		}
	} else {
		log.Println("Admin User already exists, skipping creation.")
	}

	log.Println("Seeding process completed.")
}

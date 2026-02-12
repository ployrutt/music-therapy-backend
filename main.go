package main

import (
	"log"

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
	// โหลด .env file (ถ้ามี) - ใน Docker จะใช้ env vars จาก Dokploy แทน
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	cfg := config.GetDBConfig()

	// 1. เชื่อมต่อ Database
	gormDB, err := db.InitDB(
		cfg,
		&models.User{},
		&models.Role{},
		&models.PermissionGroup{},
		&models.Permission{},
		&models.Activity{},
		&models.ActivityGoal{},
		&models.ActivitySubGoal{},
		&models.ActivityMainCategory{},
		&models.ActivitySubCategory{},

		&models.UserFavorite{},
		&models.UserReadHistory{},
	)

	if err != nil {
		log.Fatalf("Application startup failed: %v", err)
	}

	// 2. ล้างตารางเดิมทิ้ง (สำคัญมากสำหรับ PostgreSQL เพื่อไม่ให้ติด FK Error)
	// log.Println("Resetting database tables...")
	// gormDB.Exec(`DROP TABLE IF EXISTS
	// 	role_permissions, role_permission_groups,
	// 	activity_selected_sub_goals, activity_selected_sub_categories,
	// 	activity_sub_goals, activity_sub_categories,
	// 	activity_goals, activity_main_categories,
	// 	activities, users, permissions, permission_groups, roles CASCADE`)

	// 3. สร้างตารางใหม่ทั้งหมด
	log.Println("Migrating database structure...")

	err = gormDB.AutoMigrate(
		&models.Role{},
		&models.PermissionGroup{},
		&models.Permission{},
		&models.User{},
		&models.Activity{},
		&models.ActivityGoal{},
		&models.ActivitySubGoal{},
		&models.ActivityMainCategory{},
		&models.ActivitySubCategory{},
	)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// 4. รันข้อมูลเริ่มต้น
	runDatabaseSeeds(gormDB)

	log.Println("Database connection and migration successful.")

	// 5. เริ่มต้น Server
	r := router.SetupRouter(gormDB)
	port := ":8080"
	log.Printf("Starting HTTP server on port %s", port)

	sqlDB, _ := gormDB.DB()
	defer sqlDB.Close()

	if err := r.Run(port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
func runDatabaseSeeds(gormDB *gorm.DB) {
	log.Println("Starting database seeding...")

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

	// var adminCount int64

	// gormDB.Model(&models.User{}).Where("email = ?", "admin@example.com").Count(&adminCount)

	// if adminCount == 0 {

	// 	hashedPass, _ := hashPassword("password1234")
	// 	adminUser := models.User{
	// 		FirstName:   "test",
	// 		LastName:    "admin",
	// 		DateOfBirth: time.Date(1998, time.March, 1, 0, 0, 0, 0, time.UTC),
	// 		Email:       "admin@example.com",
	// 		Password:    hashedPass,
	// 		PhoneNumber: "0123456789",
	// 		Profile:     "image",
	// 		RoleID:      1,
	// 	}
	// 	if err := gormDB.Create(&adminUser).Error; err != nil {
	// 		log.Printf("Could not create initial Admin: %v", err)
	// 	} else {
	// 		log.Println("Admin User created successfully.")
	// 	}
	// } else {
	// 	log.Println("Admin User already exists, skipping...")
	// }

	log.Println("Seeding process completed.")
}

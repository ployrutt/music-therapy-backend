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

	if err := godotenv.Load(); err != nil {
		log.Println("Dokploy Environment detected: Using system variables")
	}

	cfg := config.GetDBConfig()

	gormDB, err := db.InitDB(
		cfg,
		&models.Role{},
		&models.PermissionGroup{},
		&models.Permission{},
		&models.User{},
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

	runDatabaseSeeds(gormDB)

	log.Println("Database connection and migration successful.")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := router.SetupRouter(gormDB)
	log.Printf("Starting HTTP server on port %s in %s mode", port, os.Getenv("GIN_MODE"))

	if err := r.Run(":" + port); err != nil {
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

	var adminCount int64

	targetEmail := "TestAdmin@example.com"
	gormDB.Model(&models.User{}).Where("email = ?", targetEmail).Count(&adminCount)

	if adminCount == 0 {
		hashedPass, _ := hashPassword("12052546")
		adminUser := models.User{
			FirstName:   "test",
			LastName:    "admin",
			DateOfBirth: time.Date(1999, time.March, 1, 0, 0, 0, 0, time.UTC),
			Email:       targetEmail,
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

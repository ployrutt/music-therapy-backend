package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"project-backend/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(cfg *config.DBConfig, models ...interface{}) (*gorm.DB, error) {

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	err = db.AutoMigrate(models...)
	if err != nil {
		return nil, fmt.Errorf("failed to auto migrate tables: %w", err)
	}

	fmt.Println("Database migration completed successfully!")
	return db, nil
}

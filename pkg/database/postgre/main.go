package database

import (
	"fmt"
	"golang-skeleton/internal/models"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

func NewPostgreConn(databaseURL string) (*gorm.DB, error) {
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	db, err := gorm.Open(postgres.Open(databaseURL), config)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to db: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("Failed to get database instance: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(1 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("Failed to ping database: %w", err)
	}

	log.Println("Connected to database")
	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	log.Println("Migrating database...")

	err := db.AutoMigrate(&models.User{}, &models.Role{})
	if err != nil {
		return fmt.Errorf("Failed to migrate database: %w", err)
	}

	log.Println("Database migrated")
	return nil
}

func SeedDatabase(db *gorm.DB) error {
	log.Println("Seeding database...")

	roles := []models.Role{
		{
			Name: "Admin", Code: "admin",
		},
		{
			Name: "User", Code: "user",
		},
	}

	for _, role := range roles {
		if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&role).Error; err != nil {
			return err
		}
	}

	log.Println("Database seeded")
	return nil
}

func CreateIndexes(db *gorm.DB) error {
	log.Println("Creating database indexes...")

	if err := db.Exec(`
        CREATE INDEX IF NOT EXISTS idx_users_email_idx ON users(email)
    `).Error; err != nil {
		return err
	}

	if err := db.Exec(`
        CREATE INDEX IF NOT EXISTS idx_users_role_idx ON users(role_id)
    `).Error; err != nil {
		return err
	}

	return nil
}

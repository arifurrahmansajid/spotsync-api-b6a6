package config

import (
	"fmt"
	"log"
	"os"
	"spotsync/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConnectDB initializes and returns a GORM PostgreSQL connection.
// Connection details are read from environment variables.
func ConnectDB() *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		host := getEnv("DB_HOST", "localhost")
		user := getEnv("DB_USER", "postgres")
		password := getEnv("DB_PASSWORD", "postgres")
		dbname := getEnv("DB_NAME", "spotsync")
		port := getEnv("DB_PORT", "5432")
		sslmode := getEnv("DB_SSLMODE", "disable")
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
			host, user, password, dbname, port, sslmode,
		)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate tables in dependency order
	if err := db.AutoMigrate(
		&models.User{},
		&models.ParkingZone{},
		&models.Reservation{},
	); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	log.Println("Database connected and migrated successfully")
	return db
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
